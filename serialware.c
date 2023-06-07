#include "stm32l4xx_hal.h"
#include <serialware.h>

#define gpio_low(pin)    HAL_GPIO_WritePin(pin##_GPIO_Port, pin##_Pin, GPIO_PIN_RESET)
#define gpio_high(pin)    HAL_GPIO_WritePin(pin##_GPIO_Port, pin##_Pin, GPIO_PIN_SET)
#define gpio_ishigh(pin)    (HAL_GPIO_ReadPin(pin##_GPIO_Port, pin##_Pin) == GPIO_PIN_SET)
#define gpio_toggle(pin)    HAL_GPIO_TogglePin(pin##_GPIO_Port, pin##_Pin)
#define select_rpi() HAL_GPIO_WritePin(GPIOC, SPI3_MUX_S_Pin, GPIO_PIN_SET)
#define select_leds() HAL_GPIO_WritePin(GPIOC, SPI3_MUX_S_Pin, GPIO_PIN_RESET)
#define enable_mux_out() HAL_GPIO_WritePin(GPIOC, SPI3_MUX_OE_Pin, GPIO_PIN_RESET)
#define disable_mux_out() HAL_GPIO_WritePin(GPIOC, SPI3_MUX_OE_Pin, GPIO_PIN_SET)
#define status_led_high() HAL_GPIO_WritePin(GPIOC, LED5_Pin, GPIO_PIN_SET)
#define status_led_low() HAL_GPIO_WritePin(GPIOC, LED5_Pin, GPIO_PIN_RESET)
#define status_led_toggle() HAL_GPIO_TogglePin(GPIOC, LED5_Pin)

extern UART_HandleTypeDef huart1;
extern SPI_HandleTypeDef hspi3;

static void msec_delay(int n) {
    HAL_Delay(n);
}

static void spi_detach(void) {
    HAL_SPI_MspDeInit(&hspi3);
    HAL_GPIO_DeInit(ICE40_SPI_CS_GPIO_Port, ICE40_SPI_CS_Pin);
}

static void spi_reattach(void) {
    GPIO_InitTypeDef g;

    HAL_SPI_MspInit(&hspi3);
    g.Pin = ICE40_SPI_CS_Pin;
    g.Mode = GPIO_MODE_OUTPUT_PP;
    g.Pull = GPIO_NOPULL;
    g.Speed = GPIO_SPEED_FREQ_LOW;
    HAL_GPIO_Init(ICE40_SPI_CS_GPIO_Port, &g);
}


uint8_t in_len;
uint8_t in_buf[255];

uint8_t out_len;
uint8_t out_buf[255];

static uint8_t spi_write(uint8_t *p, uint32_t len) {
    int ret;
    uint16_t n;

    ret = HAL_OK;
    n = 0x8000;
    while (len > 0) {
        if (len < n)
            n = len;
        ret = HAL_SPI_Transmit(&hspi3, p, n, HAL_MAX_DELAY);
        if (ret != HAL_OK)
            return ret;
        len -= n;
        p += n;
    }
    return ret;
}


const uint8_t CRC7_POLY = 0x91;

uint8_t getCRC(const uint8_t message[], uint8_t length) {
    uint8_t i, j, crc = 0;

    for (i = 0; i < length; i++) {
        crc ^= message[i];
        for (j = 0; j < 8; j++) {
            if (crc & 1)
                crc ^= CRC7_POLY;
            crc >>= 1;
        }
    }
    return crc;
}


uint8_t readPacket() {
    while (1) {
        uint8_t m;
        if (HAL_UART_Receive(&huart1, &m, 1, 5000) != HAL_OK) {
            continue;
        }

        if (m != 0x5c) {
            continue;
        }

        uint8_t t;

        if (HAL_UART_Receive(&huart1, &t, 1, 5000) != HAL_OK) {
            continue;
        }

        if (HAL_UART_Receive(&huart1, &in_len, 1, 500) != HAL_OK) {
            continue;
        }

        if (in_len > 0) {
            if (HAL_UART_Receive(&huart1, in_buf, in_len, 5000) != HAL_OK) {
                continue;
            }
        }

        uint8_t calcCRC;
        uint8_t expectCRC;

        if (HAL_UART_Receive(&huart1, &calcCRC, 1, 500) != HAL_OK) {
            continue;
        }

        expectCRC = getCRC(in_buf, in_len);

        if (expectCRC != calcCRC) {
            uint8_t b = 0x5d;
            HAL_UART_Transmit(&huart1, (unsigned char *) &b, 1, 500);

            continue;
        }

        uint8_t b = 0x5e;
        HAL_UART_Transmit(&huart1, (unsigned char *) &b, 1, 500);

        return t;
    }
}

void writePacket(uint8_t id) {
    uint8_t b = 0x5c;
    HAL_UART_Transmit(&huart1, (unsigned char *) &b, 1, 500);
    HAL_UART_Transmit(&huart1, (unsigned char *) &id, 1, 500);
    HAL_UART_Transmit(&huart1, (unsigned char *) &out_len, 1, 500);
    HAL_UART_Transmit(&huart1, (unsigned char *) out_buf, out_len, 5000);

    uint8_t calcCRC;
    calcCRC = getCRC(out_buf, out_len);

    HAL_UART_Transmit(&huart1, (unsigned char *) &calcCRC, 1, 500);
}

static uint8_t ice40_reset(void) {
    int timeout;

    gpio_low(ICE40_CRST);
    gpio_low(ICE40_SPI_CS);
    msec_delay(1);
    gpio_high(ICE40_CRST);
    timeout = 100;
    while (gpio_ishigh(ICE40_CDONE)) {
        if (--timeout == 0)
            return 1;
    }
    msec_delay(2);
    return 0;
}

static uint8_t ice40_configdone(void) {
    uint8_t b = 0;

    for (int timeout = 100; !gpio_ishigh(ICE40_CDONE); timeout--) {
        if (timeout == 0) {
            return 1;
        }

        spi_write(&b, 1);
    }

    for (int i = 0; i < 7; i++)
        spi_write(&b, 1);

    return 0;
}


_Noreturn void run(void) {
    select_leds();

    uint8_t count = 0;

    while (1) {
        uint8_t id;
        id = readPacket();

        uint8_t res;

        switch (id) {
            case 5:
                out_len = 1;
                out_buf[0] = 123;
                writePacket(id);
                break;

            case 10:
                status_led_high();
                disable_mux_out();
                spi_reattach();

                res = ice40_reset();

                out_len = 1;
                out_buf[0] = res;
                writePacket(id);

                count = 0;

                break;

            case 11:
                count++;

                if (count == 10) {
                    count = 0;
                    status_led_toggle();
                }

                res = spi_write(in_buf, in_len);

                out_len = 1;
                out_buf[0] = res;
                writePacket(id);
                break;

            case 12:
                res = ice40_configdone();

                spi_detach();
                enable_mux_out();
                status_led_low();

                out_len = 1;
                out_buf[0] = res;
                writePacket(id);
                break;

            default:
                break;
        }
    }
}
