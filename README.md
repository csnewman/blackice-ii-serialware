# blackice-ii-serialware

Alternative firmware for the BlackIce II to allow uploads over serial.

This firmware is hacky and should be used at your own risk.

### Build Guide

1. `mkdir build`
2. `cd build`
3. `cmake ..`
4. `build`

`serialware.bin` should now exist in the build directory and be ready for flashing.

### Flashing Guide

1. Connect SWD device (e.g. Raspberry PI or STLink)
2. Create `openocd.cfg`, replacing the interface file and adapter speed as necessary
    ```
    source /usr/local/share/openocd/scripts/interface/raspberrypi-native.cfg
    transport select swd
    source /usr/local/share/openocd/scripts/target/stm32l4x.cfg
    reset_config  srst_nogate
    
    adapter_nsrst_delay 100
    adapter_nsrst_assert_width 100
    adapter speed 5
    
    init
    targets
    ```
3. Start openocd
    ```
    sudo openocd
    ```
4. Connect to openocd
    ```
    telnet localhost 4444
    ```
5. Flash
    ```
    halt
    flash write_image erase serialware.bin 0x08000000
    ```
6. Reset device

### Thanks

The firmware is based upon the
original [BlackIce II firmware](https://github.com/mystorm-org/BlackIce-II/tree/master/firmware).
