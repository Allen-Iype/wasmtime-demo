extern "C" {
    fn load_input(pointer: *mut u8);
    fn dump_output(pointer: *const u8, user_rating: f32 , rating_count: f32, length: usize);
}

//did,rating,count

#[no_mangle]
pub extern "C" fn handler(did_length: usize , rating_length: usize , rating_count_length: usize, user_rating_length: usize) {
    // load input data
    let mut input = Vec::with_capacity(did_length + rating_length + rating_count_length + user_rating_length);
    
    unsafe {
        load_input(input.as_mut_ptr());
        input.set_len(did_length + rating_length + rating_count_length + user_rating_length);
    }


    let (did, b1_rest) = input.split_at(did_length);
    let (rating, count_rating) = b1_rest.split_at(rating_length);
    let (rating_count, user_rating) = count_rating.split_at(rating_count_length);
   
    
    let mut current_rating = f32::from_ne_bytes(rating[0..rating_length].try_into().unwrap());
    let mut total_count = f32::from_ne_bytes(rating_count[0..rating_count_length].try_into().unwrap());
    let latest_user_rating = f32::from_ne_bytes(user_rating[0..user_rating_length].try_into().unwrap());
    
    total_count += 1.00;
    current_rating = (current_rating + latest_user_rating)/(total_count + 1.00);

    // dump output data
    unsafe {
        dump_output(did.as_ptr() , current_rating, total_count, did.len());

    }
}
