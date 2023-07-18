extern "C" {
    fn load_input(pointer: *mut u8);
    fn dump_output(pointer: *const u8, user_rating: f32 , rating_count: f32, product_id_length: usize,seller_did_length: usize, current_seller_rating: f32, seller_did_pointer: *const u8);
    fn get_account_info(did: *const u8, port: u32);
}

//did,rating,count

#[no_mangle]
pub extern "C" fn handler(did_length: usize , rating_length: usize , rating_count_length: usize, user_rating_length: usize, seller_did_length: usize, seller_rating_length: usize, seller_product_count_length: usize) {
    // load input data
    let mut input = Vec::with_capacity(did_length + rating_length + rating_count_length + user_rating_length + seller_did_length + seller_rating_length + seller_product_count_length);
    
    unsafe {
        load_input(input.as_mut_ptr());
        input.set_len(did_length + rating_length + rating_count_length + user_rating_length + seller_did_length + seller_rating_length + seller_product_count_length);
    }


    let (product_id, b1_rest) = input.split_at(did_length);
    let (seller_did, b2_rest) = b1_rest.split_at(seller_did_length);
    let (seller_rating, b3_rest) = b2_rest.split_at(seller_rating_length);
    let (seller_product_count,b4_rest) = b3_rest.split_at(seller_product_count_length);
    let (user_rating, b5_rest) = b4_rest.split_at(user_rating_length);
    let (rating_count, latest_rating) = b5_rest.split_at(rating_count_length);
   
    
    let mut current_rating = f32::from_ne_bytes(latest_rating[0..rating_length].try_into().unwrap());
    let mut total_count = f32::from_ne_bytes(rating_count[0..rating_count_length].try_into().unwrap());
    let latest_user_rating = f32::from_ne_bytes(user_rating[0..user_rating_length].try_into().unwrap());
  //  let mut current_seller_rating:f32 = f32::from_ne_bytes(seller_rating[0..seller_rating_length].try_into().unwrap());
  //  let total_seller_product_count = f32::from_ne_bytes(seller_product_count[0..seller_product_count_length].try_into().unwrap());
    total_count += 1.00;
    current_rating = (current_rating + latest_user_rating)/(total_count);

  //  current_seller_rating = ()/(total_seller_product_count);
    let mut seller_did_vec:Vec<u8> = Vec::with_capacity(seller_did.len());
    seller_did_vec.extend_from_slice(seller_did);

    // dump output data
    unsafe {
        dump_output(product_id.as_ptr() , current_rating, total_count, product_id.len(),seller_did_vec.len(), 1.0, seller_did.as_ptr());

    }
}


// #[no_mangle]
// pub extern "C" fn handler2(
//     did_length: usize,
//     rating_length: usize,
//     rating_count_length: usize,
//     user_rating_length: usize,
//     seller_did_length: usize,
//     seller_rating_length: usize,
//     seller_product_count_length: usize,
// ) {
//     let input_size = did_length
//         + rating_length
//         + rating_count_length
//         + user_rating_length
//         + seller_did_length
//         + seller_rating_length
//         + seller_product_count_length;

//     let input = unsafe {
//         let mut input = Vec::with_capacity(input_size);
//         load_input(input.as_mut_ptr());
//         Vec::from_raw_parts(input.as_mut_ptr(), input_size, input.capacity())
//     };

//     let (did, rest) = input.split_at(did_length);
//     let (rating, rest) = rest.split_at(rating_length);
//     let (rating_count, rest) = rest.split_at(rating_count_length);
//     let (user_rating, rest) = rest.split_at(user_rating_length);
//     let (seller_did, rest) = rest.split_at(seller_did_length);
//     let (seller_rating, seller_product_count) = rest.split_at(seller_rating_length);

//     let current_rating = f32::from_ne_bytes(rating.try_into().unwrap());
//     let total_count = f32::from_ne_bytes(rating_count.try_into().unwrap());
//     let latest_user_rating = f32::from_ne_bytes(user_rating.try_into().unwrap());
//     let current_seller_rating = f32::from_ne_bytes(seller_rating.try_into().unwrap());
//     let total_seller_product_count = f32::from_ne_bytes(seller_product_count.try_into().unwrap());

//     let new_count = total_count + 1.00;
//     let new_rating = (current_rating * total_count + latest_user_rating) / new_count;

//     // dump output data
//     unsafe {
//         dump_output(did.as_ptr(), new_rating, new_count, did_length);
//     }
// }
