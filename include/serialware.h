#ifndef __SERIALWARE_H
#define __SERIALWARE_H

_Noreturn void run(void);

#define LED5_Pin GPIO_PIN_13
#define LED5_GPIO_Port GPIOC
#define SPI3_MUX_S_Pin GPIO_PIN_14
#define SPI3_MUX_S_GPIO_Port GPIOC
#define SPI3_MUX_OE_Pin GPIO_PIN_15
#define SPI3_MUX_OE_GPIO_Port GPIOC
#define AN0_Pin GPIO_PIN_0
#define AN0_GPIO_Port GPIOC
#define AN1_Pin GPIO_PIN_1
#define AN1_GPIO_Port GPIOC
#define AN2_Pin GPIO_PIN_2
#define AN2_GPIO_Port GPIOC
#define AN3_Pin GPIO_PIN_3
#define AN3_GPIO_Port GPIOC
#define AN4_Pin GPIO_PIN_0
#define AN4_GPIO_Port GPIOA
#define AN5_Pin GPIO_PIN_1
#define AN5_GPIO_Port GPIOA
#define QSPI_CS_Pin GPIO_PIN_2
#define QSPI_CS_GPIO_Port GPIOA
#define QSPI_CLK_Pin GPIO_PIN_3
#define QSPI_CLK_GPIO_Port GPIOA
#define DIG9_Pin GPIO_PIN_4
#define DIG9_GPIO_Port GPIOA
#define DIG5_Pin GPIO_PIN_5
#define DIG5_GPIO_Port GPIOA
#define QSPI_D3_Pin GPIO_PIN_6
#define QSPI_D3_GPIO_Port GPIOA
#define QSPI_D2_Pin GPIO_PIN_7
#define QSPI_D2_GPIO_Port GPIOA
#define DIG0_Pin GPIO_PIN_4
#define DIG0_GPIO_Port GPIOC
#define DIG1_Pin GPIO_PIN_5
#define DIG1_GPIO_Port GPIOC
#define QSPI_D1_Pin GPIO_PIN_0
#define QSPI_D1_GPIO_Port GPIOB
#define QSPI_D0_Pin GPIO_PIN_1
#define QSPI_D0_GPIO_Port GPIOB
#define DIG2_Pin GPIO_PIN_2
#define DIG2_GPIO_Port GPIOB
#define DIG3_Pin GPIO_PIN_10
#define DIG3_GPIO_Port GPIOB
#define DIG4_Pin GPIO_PIN_11
#define DIG4_GPIO_Port GPIOB
#define DIG10_Pin GPIO_PIN_12
#define DIG10_GPIO_Port GPIOB
#define DIG13_Pin GPIO_PIN_13
#define DIG13_GPIO_Port GPIOB
#define DIG12_Pin GPIO_PIN_14
#define DIG12_GPIO_Port GPIOB
#define DIG11_Pin GPIO_PIN_15
#define DIG11_GPIO_Port GPIOB
#define DIG7_Pin GPIO_PIN_6
#define DIG7_GPIO_Port GPIOC
#define DIG8_Pin GPIO_PIN_7
#define DIG8_GPIO_Port GPIOC
#define B1_Pin GPIO_PIN_8
#define B1_GPIO_Port GPIOC
#define B2_Pin GPIO_PIN_9
#define B2_GPIO_Port GPIOC
#define DIG6_Pin GPIO_PIN_8
#define DIG6_GPIO_Port GPIOA
#define ICE40_SPI_CS_Pin GPIO_PIN_15
#define ICE40_SPI_CS_GPIO_Port GPIOA
#define DIG17_Pin GPIO_PIN_10
#define DIG17_GPIO_Port GPIOC
#define DIG18_Pin GPIO_PIN_11
#define DIG18_GPIO_Port GPIOC
#define DIG19_Pin GPIO_PIN_12
#define DIG19_GPIO_Port GPIOC
#define DIG16_Pin GPIO_PIN_2
#define DIG16_GPIO_Port GPIOD
#define ICE40_CRST_Pin GPIO_PIN_6
#define ICE40_CRST_GPIO_Port GPIOB
#define ICE40_CDONE_Pin GPIO_PIN_7
#define ICE40_CDONE_GPIO_Port GPIOB



#define QSPICSN_PIN GPIO_PIN_2
#define QSPICSN_PORT GPIOA
#define QSPICK_PIN GPIO_PIN_3
#define QSPICK_PORT GPIOA
#define QSPIDQ0_PIN GPIO_PIN_1
#define QSPIDQ0_PORT GPIOB
#define QSPIDQ1_PIN GPIO_PIN_0
#define QSPIDQ1_PORT GPIOB
#define QSPIDQ2_PIN GPIO_PIN_7
#define QSPIDQ2_PORT GPIOA
#define QSPIDQ3_PIN GPIO_PIN_6
#define QSPIDQ3_PORT GPIOA


#define QCLK_Pin QSPICK_PIN
#define QCLK_GPIO_Port QSPICK_PORT
#define QDIR_Pin QSPICSN_PIN
#define QDIR_GPIO_Port QSPICSN_PORT
#define QCTS_Pin QSPIDQ0_PIN
#define QCTS_GPIO_Port QSPIDQ0_PORT
#define QRTS_Pin QSPIDQ1_PIN
#define QRTS_GPIO_Port QSPIDQ1_PORT
#define QIN_Pin QSPIDQ2_PIN
#define QIN_GPIO_Port QSPIDQ2_PORT
#define QOUT_Pin QSPIDQ3_PIN
#define QOUT_GPIO_Port QSPIDQ3_PORT


#define QSLEEP 100

#endif