cmake_minimum_required(VERSION 3.25)
project(serialware)

enable_language(C ASM)
set(CMAKE_C_STANDARD 11)
set(CMAKE_C_STANDARD_REQUIRED ON)
set(CMAKE_C_EXTENSIONS OFF)
set(CMAKE_SYSTEM_NAME Generic)
set(CMAKE_SYSTEM_PROCESSOR ARM)
set(UTIL_SEARCH_CMD which)
set(TOOLCHAIN_PREFIX arm-none-eabi-)

execute_process(
        COMMAND ${UTIL_SEARCH_CMD} ${TOOLCHAIN_PREFIX}gcc
        OUTPUT_VARIABLE BINUTILS_PATH
        OUTPUT_STRIP_TRAILING_WHITESPACE
)

get_filename_component(ARM_TOOLCHAIN_DIR ${BINUTILS_PATH} DIRECTORY)
set(CMAKE_TRY_COMPILE_TARGET_TYPE STATIC_LIBRARY)

set(CMAKE_C_COMPILER ${TOOLCHAIN_PREFIX}gcc)
set(CMAKE_ASM_COMPILER ${CMAKE_C_COMPILER})
set(CMAKE_CXX_COMPILER ${TOOLCHAIN_PREFIX}g++)

set(CMAKE_OBJCOPY ${ARM_TOOLCHAIN_DIR}/${TOOLCHAIN_PREFIX}objcopy CACHE INTERNAL "objcopy tool")
set(CMAKE_SIZE_UTIL ${ARM_TOOLCHAIN_DIR}/${TOOLCHAIN_PREFIX}size CACHE INTERNAL "size tool")

set(CMAKE_FIND_ROOT_PATH ${BINUTILS_PATH})
set(CMAKE_FIND_ROOT_PATH_MODE_PROGRAM NEVER)
set(CMAKE_FIND_ROOT_PATH_MODE_LIBRARY ONLY)
set(CMAKE_FIND_ROOT_PATH_MODE_INCLUDE ONLY)


set(EXECUTABLE ${PROJECT_NAME}.out)

set(STM32CUBEMX_GENERATED_FILES
        main.c
        serialware.c
        stm32l4xx_hal_msp.c
        stm32l4xx_it.c

        STM32CubeL4/Drivers/CMSIS/Device/ST/STM32L4xx/Source/Templates/gcc/startup_stm32l433xx.s
        STM32CubeL4/Drivers/CMSIS/Device/ST/STM32L4xx/Source/Templates/system_stm32l4xx.c
        STM32CubeL4/Drivers/STM32L4xx_HAL_Driver/Src/stm32l4xx_hal.c
        STM32CubeL4/Drivers/STM32L4xx_HAL_Driver/Src/stm32l4xx_hal_cortex.c
        STM32CubeL4/Drivers/STM32L4xx_HAL_Driver/Src/stm32l4xx_hal_dma.c
        STM32CubeL4/Drivers/STM32L4xx_HAL_Driver/Src/stm32l4xx_hal_gpio.c
        STM32CubeL4/Drivers/STM32L4xx_HAL_Driver/Src/stm32l4xx_hal_i2c.c
        STM32CubeL4/Drivers/STM32L4xx_HAL_Driver/Src/stm32l4xx_hal_i2c_ex.c
        STM32CubeL4/Drivers/STM32L4xx_HAL_Driver/Src/stm32l4xx_hal_pwr_ex.c
        STM32CubeL4/Drivers/STM32L4xx_HAL_Driver/Src/stm32l4xx_hal_qspi.c
        STM32CubeL4/Drivers/STM32L4xx_HAL_Driver/Src/stm32l4xx_hal_rcc.c
        STM32CubeL4/Drivers/STM32L4xx_HAL_Driver/Src/stm32l4xx_hal_rcc_ex.c
        STM32CubeL4/Drivers/STM32L4xx_HAL_Driver/Src/stm32l4xx_hal_spi.c
        STM32CubeL4/Drivers/STM32L4xx_HAL_Driver/Src/stm32l4xx_hal_uart.c
        STM32CubeL4/Drivers/STM32L4xx_HAL_Driver/Src/stm32l4xx_hal_uart_ex.c
        )
add_executable(${EXECUTABLE} ${STM32CUBEMX_GENERATED_FILES})

target_include_directories(${EXECUTABLE} PRIVATE
        include
        STM32CubeL4/Drivers/STM32L4xx_HAL_Driver/Inc
        STM32CubeL4/Drivers/CMSIS/Device/ST/STM32L4xx/Include
        STM32CubeL4/Drivers/CMSIS/Include
        )

target_compile_definitions(${EXECUTABLE} PRIVATE -DSTM32L433xx)
target_compile_options(${EXECUTABLE} PRIVATE
        -mcpu=cortex-m3
        -mthumb
        -fdata-sections
        -ffunction-sections
        -Wall
        -Werror
        -O0
        )
#        $<$<CONFIG:Debug>:-Og>

target_link_options(${EXECUTABLE} PRIVATE
        -T${CMAKE_SOURCE_DIR}/blackice_flash.ld
        -mcpu=cortex-m3
        -mthumb
        -lc
        -lnosys
        -Wl,--gc-sections
        )

add_custom_command(TARGET ${EXECUTABLE}
        POST_BUILD
        COMMAND arm-none-eabi-size ${EXECUTABLE}
        )

add_custom_command(TARGET ${EXECUTABLE}
        POST_BUILD
        COMMAND arm-none-eabi-objcopy -O ihex ${EXECUTABLE} ${PROJECT_NAME}.hex
        COMMAND arm-none-eabi-objcopy -O binary ${EXECUTABLE} ${PROJECT_NAME}.bin
        )
