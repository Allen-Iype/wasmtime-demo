extern "C" {
    fn load_input(pointer: *mut u8);
    fn dump_output(pointer: *const u8, uservote: u32 , redvote: u32 , bluevote: u32 , length: usize);
}

#[no_mangle]
pub extern "C" fn handler(input_length: usize , b1_length: usize , redcount: usize , bluecount: usize) {
    // load input data
    let mut input = Vec::with_capacity(input_length + b1_length + redcount + bluecount);
    
    unsafe {
        load_input(input.as_mut_ptr());
        input.set_len(input_length + b1_length + redcount + bluecount);
    }


    let (buf, b1_rest) = input.split_at(input_length);
    let (b1, red_blue) = b1_rest.split_at(b1_length);
    let (red, blue) = red_blue.split_at(redcount);

    // process app data
    let output = buf.to_ascii_uppercase();
    println!("Did ouput{:?}",output);

    let uservote = u32::from_ne_bytes(b1[0..b1_length].try_into().unwrap());
    let mut redvote = u32::from_ne_bytes(red[0..redcount].try_into().unwrap());
    let mut bluevote = u32::from_ne_bytes(blue[0..bluecount].try_into().unwrap());
    
    if uservote == 1 {
    redvote += 1;
    }
    else {
        bluevote += 1
    }

    // dump output data
    unsafe {
        dump_output(output.as_ptr() , uservote, redvote , bluevote , output.len());

    }
}
