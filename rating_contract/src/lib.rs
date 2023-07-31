use serde::{Serialize, Deserialize};
use serde_cbor;

#[derive(Debug, Serialize, Deserialize)]
struct  ProductReview {
    #[serde(rename = "product_id")]
	product_id:   String,
    #[serde(rename = "rating")]
	rating:      f32,
    #[serde(rename = "rating_count")]
	rating_count: f32,
    #[serde(rename = "seller_did")]
	seller_did:   String,
}

#[derive(Debug, Serialize, Deserialize)]
struct SellerReview {
    #[serde(rename = "did")]
	did:          String,
    #[serde(rename = "seller_rating")]
	seller_rating: f32,
    #[serde(rename = "product_count")]
	product_count: f32,
}


extern "C" {
    fn load_input(pointer: *mut u8);
    fn dump_output(pointer: *const u8, product_review_len: usize , seller_review_len: usize);
    fn get_account_info();
    fn initiate_transfer();
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
    let product_review: ProductReview = serde_cbor::from_slice(product_state).expect("Failed to decode CBOR data");
    let seller_review: SellerReview = serde_cbor::from_slice(seller_state).expect("Failed to decode CBOR data");
    //append ProductReview and SellerReview
    let product_id = product_review.product_id;
    let current_rating = product_review.rating;
    let rating_count = product_review.rating_count;
    let seller_did_product = product_review.seller_did;

    let seller_did = seller_review.did;
    let current_seller_rating = seller_review.seller_rating;
    let seller_product_count = seller_review.product_count;

    let new_count = rating_count + 1.00;
    let new_rating = ((current_rating * rating_count) + rating) / new_count;
    // if new_count == 10.0 && new_rating>= 3.00 {
    //     // initiate transfer
    //     unsafe {
    //         initiate_transfer();
    //     }
    // }
    let new_seller_rating = (current_seller_rating + new_rating) / seller_product_count;

    let product_review_test = ProductReview{ product_id: product_id, rating: new_rating, rating_count: new_count, seller_did: seller_did_product };
    let seller_review_test = SellerReview{ did: seller_did, seller_rating: new_seller_rating, product_count: seller_product_count };
    let cbor_product_review:Vec<u8> = serde_cbor::to_vec(&product_review_test).expect("Failed to serialize to CBOR");
    let cbor_seller_review:Vec<u8> = serde_cbor::to_vec(&seller_review_test).expect("Failed to serialize to CBOR");
    let latest_product_len = cbor_product_review.len();
    let latest_seller_len = cbor_seller_review.len();
  //  current_seller_rating = ()/(total_seller_product_count);
    // append two vectors
    let combined_vec = [cbor_product_review, cbor_seller_review].concat();
    unsafe {
        dump_output(combined_vec.as_ptr() , latest_product_len,latest_seller_len);
    }
    if new_count == 10.0 && new_rating>= 3.00 {
        // initiate transfer
        unsafe {
         //   dump_output(combined_vec.as_ptr() , latest_product_len,latest_seller_len);
            get_account_info();
            initiate_transfer();
        }
    }

    // dump output data
    // unsafe {
    //    dump_output(combined_vec.as_ptr() , latest_product_len,latest_seller_len);
    //    initiate_transfer();
    //    get_account_info();
      

    // }
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