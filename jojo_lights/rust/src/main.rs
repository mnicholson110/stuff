#![deny(unsafe_code)]
#![no_main]
#![no_std]

use cortex_m_rt::entry;
use panic_halt as _;
use microbit::{
    board::Board,
    display::blocking::Display,
    hal::Timer,
};

const PIXELS: [(usize,usize); 25] = [
    (0,0), (0,1), (0,2), (0,3), (0,4),
    (1,0), (1,1), (1,2), (1,3), (1,4),
    (2,0), (2,1), (2,2), (2,3), (2,4),
    (3,0), (3,1), (3,2), (3,3), (3,4),
    (4,0), (4,1), (4,2), (4,3), (4,4),
];

#[entry]
fn main() -> ! {

    let board = Board::take().unwrap();
    let mut timer = Timer::new(board.TIMER0);
    let mut display = Display::new(board.display_pins);
    let mut leds = [
        [0,0,0,0,0],
        [0,0,0,0,0],
        [0,0,0,0,0],
        [0,0,0,0,0],
        [0,0,0,0,0],
    ];

    let mut last_led = (0,0);

    loop {
        for current_led in PIXELS.iter() {
            leds[last_led.0][last_led.1] = 0;
            leds[current_led.0][current_led.1] = 9;
            display.show(&mut timer, leds, 300);
            last_led = *current_led;
        }
    }
}


