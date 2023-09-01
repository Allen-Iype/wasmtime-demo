
extern "C" {
    fn load_input(pointer: *mut u8);
    fn dump_output(pointer: *const u8, redvote: u32 , bluevote: u32 , port_length: usize, hash_length: usize);
}

fn dumy(){
    println!("")
}

fn dummy2(){
    print!("")
}

fn dummy3() {

}

fn dummy4() {
    
}

fn dummy5() {

}

fn dummy6() {

}

fn dummy7() {

}

fn dummy8() {

}

fn dummy9() {

}

fn dummy10() {

}

fn dummy11() {

}

fn dummyfjdgd(){
    println!();
}

fn dummyLatest() {

}

#[no_mangle]
pub extern "C" fn handler(input_vote_length: usize , red_length: usize , blue_length: usize, port_length: usize, hash_length: usize) {
    // load input data
    let mut input = Vec::with_capacity(input_vote_length + red_length + blue_length + port_length + hash_length);
    let mut output_vec:Vec<u8> = Vec::new();
    unsafe {
        load_input(input.as_mut_ptr());
        input.set_len(input_vote_length + red_length + blue_length + port_length + hash_length);
    
    }


    let (input_vote, b1_rest) = input.split_at(input_vote_length);
    let (red_count, blue_port_hash) = b1_rest.split_at(red_length);
    let (blue_count, port_hash) = blue_port_hash.split_at(blue_length);
    let (port_byte, hash_byte) = port_hash.split_at(port_length);


    if let Ok(user_vote) = std::str::from_utf8(&input_vote) {
        let mut red_vote = u32::from_ne_bytes(red_count[0..red_length].try_into().unwrap());
    let mut blue_vote = u32::from_ne_bytes(blue_count[0..blue_length].try_into().unwrap());
   
    if user_vote == "Red" {
        red_vote += 1;
    } else {
        blue_vote += 1;
    }

    output_vec.extend_from_slice(port_byte);
    output_vec.extend_from_slice(hash_byte);
    // dump output data
    unsafe {
        dump_output(output_vec.as_ptr(), red_vote , blue_vote , port_byte.len(),hash_byte.len());

    }
    } else {
        println!("Invalid UTF-8 sequence");
        dummy2();
        dumy() ;
        dummy3();
        dummy4();
        dummy5();
        dummy6();
        dummy10();
        dummy8();
        dummy9();
        dummy11();

        dummyLatest();
    }
}


//////////////////////////////////////////////////////
/////////////////////////////////////////////////////////
// 
//