
extern "C" {
    fn load_input(pointer: *mut u8);
    fn dump_output(pointer: *const u32, redvote: u32 , bluevote: u32 , length: usize);
}

fn dummy1(){
    println!("dummy1");
}

fn dummy2(){
    println!("dummy2");
}

fn dummy3(){
    println!("dummy3");
}

#[no_mangle]
pub extern "C" fn handler(input_vote_length: usize , red_length: usize , blue_length: usize) {
    // load input data
    let mut input = Vec::with_capacity(input_vote_length + red_length + blue_length);
    let mut output_vec:Vec<u32> = Vec::new();
    unsafe {
        load_input(input.as_mut_ptr());
        input.set_len(input_vote_length + red_length + blue_length);
    
    }


    let (input_vote, b1_rest) = input.split_at(input_vote_length);
    let (red_count, blue_count) = b1_rest.split_at(red_length);


    if let Ok(user_vote) = std::str::from_utf8(&input_vote) {
        let mut red_vote = u32::from_ne_bytes(red_count[0..red_length].try_into().unwrap());
    let mut blue_vote = u32::from_ne_bytes(blue_count[0..blue_length].try_into().unwrap());
    if user_vote == "Red" {
        red_vote += 1;
    } else {
        blue_vote += 1;
    }

    output_vec.push(red_vote);
    output_vec.push(blue_vote);

    // dump output data
    unsafe {
        dump_output(output_vec.as_ptr(), red_vote , blue_vote , output_vec.len());

    }
    } else {
        println!("Invalid UTF-8 sequence");
        dummy1();
    }
}
