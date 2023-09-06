
extern "C" {
    fn load_input(pointer: *mut u8);
    fn dump_output(pointer: *const u8, redvote: u32 , bluevote: u32 , block_no: u32, port_length: usize, hash_length: usize, block_id_length: usize);
}

#[no_mangle]
pub extern "C" fn handler(input_vote_length: usize , red_length: usize , blue_length: usize, port_length: usize, hash_length: usize, block_id_length: usize, block_no_length: usize ) {
    // load input data
    let mut input = Vec::with_capacity(input_vote_length + red_length + blue_length + port_length + hash_length + block_id_length + block_no_length);
    let mut output_vec:Vec<u8> = Vec::new();
    unsafe {
        load_input(input.as_mut_ptr());
        input.set_len(input_vote_length + red_length + blue_length + port_length + hash_length + block_id_length + block_no_length);
    
    }


    let (input_vote, b1_rest) = input.split_at(input_vote_length);
    let (red_count, blue_port_hash) = b1_rest.split_at(red_length);
    let (blue_count, port_hash) = blue_port_hash.split_at(blue_length);
    let (port_byte, hash_block_id_no) = port_hash.split_at(port_length);
    let (hash_byte,block_id_no) = hash_block_id_no.split_at(hash_length);
    let (block_id,block_no) = block_id_no.split_at(block_id_length);


    if let Ok(user_vote) = std::str::from_utf8(&input_vote) {
        let mut red_vote = u32::from_ne_bytes(red_count[0..red_length].try_into().unwrap());
    let mut blue_vote = u32::from_ne_bytes(blue_count[0..blue_length].try_into().unwrap());
    let mut block_no = u32::from_ne_bytes(block_no[0..block_no_length].try_into().unwrap());
   
    if user_vote == "Red" {
        red_vote += 1;
    } else if user_vote == "Blue" {
        blue_vote += 1;
    } else {
        println!("Invalid vote");
    }

    output_vec.extend_from_slice(port_byte);
    output_vec.extend_from_slice(hash_byte);
    output_vec.extend_from_slice(block_id);
    // dump output data
    unsafe {
        dump_output(output_vec.as_ptr(), red_vote , blue_vote , block_no, port_byte.len(),hash_byte.len(),block_id.len());

    }
    } else {
        println!("Invalid UTF-8 sequence");
    }
}


///////   ////// ///// ///// /// //