extern "C" {
    fn cast_vote(pointer: *mut u8);
    fn get_results(pointer: *const u8, length: usize);
}

#[no_mangle]
pub extern "C" fn handler(input_length: usize){
    let mut vote = Vec::with_capacity(input_length);
    unsafe {
        cast_vote(vote.as_mut_ptr());
        get_results(vote.as_ptr(), vote.len());
    }
}