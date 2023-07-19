use serde::Deserialize;

struct  ProductReview {
	ProductId:   String,
	Rating:      f32,
	RatingCount: f32,
	SellerDID:   String,
}

struct SellerReview {
	DID:          String,
	SellerRating: f32,
	ProductCount: f32,
}


extern "C" {
    fn load_input(pointer: *mut u8);
    fn dump_output(pointer: *const u8, user_rating: f32 , rating_count: f32, product_id_length: usize,seller_did_length: usize, current_seller_rating: f32, seller_did_pointer: *const u8);
    fn get_account_info(did: *const u8, port: u32);
}

//did,rating,count

#[no_mangle]
pub extern "C" fn handler(product_state_length: usize , seller_state_length: usize , rating: f32) {
    // load input data
    let mut input = Vec::with_capacity(product_state_length+seller_state_length);
    
    unsafe {
        load_input(input.as_mut_ptr());
        input.set_len(product_state_length+seller_state_length);
    }


    let (product_state, seller_state) = input.split_at(product_state_length);
   //cbor decode byte array 

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
