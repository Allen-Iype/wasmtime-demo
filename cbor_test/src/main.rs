use serde::{Serialize, Deserialize};
use serde_cbor;

#[derive(Debug, Serialize, Deserialize)]
struct  ProductReview {
	ProductId:   String,
	Rating:      f32,
	RatingCount: f32,
	SellerDID:   String,
}

#[derive(Debug, Serialize, Deserialize)]
struct SellerReview {
	DID: String,
	SellerRating: f32,
	ProductCount: f32,
}

fn main() {
    let encoded_data = vec![164, 105, 80, 114, 111, 100, 117, 99, 116, 73, 100, 107, 65, 66, 68, 70, 71, 68, 70, 71, 68, 70, 71, 102, 82, 97, 116, 105, 110, 103, 250, 63, 181, 85, 85, 107, 82, 97, 116, 105, 110, 103, 67, 111, 117, 110, 116, 250, 64, 160, 0, 0, 105, 83, 101, 108, 108, 101, 114, 68, 73, 68, 96, 163, 99, 68, 73, 68, 96, 108, 83, 101, 108, 108, 101, 114, 82, 97, 116, 105, 110, 103, 250, 63, 128, 0, 0, 108, 80, 114, 111, 100, 117, 99, 116, 67, 111, 117, 110, 116, 250, 63, 128, 0, 0];
    let (product_state, seller_state) = encoded_data.split_at(63);
    println!("product_state : {:?}",product_state);
    println!("seller_state : {:?}",seller_state);
    let product_review: ProductReview = serde_cbor::from_slice(&product_state).expect("Failed to decode CBOR data");
    println!("Product Review: {:?}", product_review);
    let seller_review: SellerReview = serde_cbor::from_slice(&seller_state).expect("Failed to decode CBOR data");
    println!("Seller Review: {:?}", seller_review);
    println!("{}", product_review.ProductId);


    let product_review_test = ProductReview{ ProductId: "Dummy".to_owned(), Rating: 1.0, RatingCount: 1.0, SellerDID: "DID".to_owned() };
    let seller_review_test = SellerReview{ DID: "DID".to_owned(), SellerRating: 2.0, ProductCount: 1.0 };
    let cbor_product_review:Vec<u8> = serde_cbor::to_vec(&product_review_test).expect("Failed to serialize to CBOR");
    let cbor_seller_review:Vec<u8> = serde_cbor::to_vec(&seller_review_test).expect("Failed to serialize to CBOR");
    let latest_product_len = cbor_product_review.len();
    let latest_seller_len = cbor_seller_review.len();
    println!("latest_product_len : {:?}",latest_product_len);
    println!("latest_seller_len : {:?}",latest_seller_len);
    println!("cbor_product_review : {:?}",cbor_product_review);
    println!("cbor_seller_review : {:?}",cbor_seller_review);
    // append two vectors
    let combined_vec = [cbor_product_review, cbor_seller_review].concat();
    println!("combined_vec : {:?}",combined_vec);
    

}


// extern "C" {
//     fn load_input(pointer: *mut u8);
//     fn dump_output(pointer: *const u8, product_review_len: usize , seller_review_len: usize);
// }

// //did,rating,count

// #[no_mangle]
// pub extern "C" fn handler(product_state_length: usize , seller_state_length: usize , rating: f32) {
//     // load input data
//     let mut input = Vec::with_capacity(product_state_length+seller_state_length);
    
//     unsafe {
//         load_input(input.as_mut_ptr());
//         input.set_len(product_state_length+seller_state_length);
//     }


//     let (product_state, seller_state) = input.split_at(product_state_length);
//     let product_review: ProductReview = serde_cbor::from_slice(product_state).expect("Failed to decode CBOR data");
//     let seller_review: SellerReview = serde_cbor::from_slice(seller_state).expect("Failed to decode CBOR data");
//     //append ProductReview and SellerReview
//     // let product_id = product_review.ProductId;
//     // let current_rating = product_review.Rating;
//     // let rating_count = product_review.RatingCount;

//     // let seller_did = seller_review.DID;
//     // let seller_rating = seller_review.SellerRating;
//     // let seller_product_count = seller_review.ProductCount;

//     // let new_count = rating_count + 1.00;
//     // let new_rating = (current_rating * rating_count + rating) / new_count;
//     // println!("new Rating is {}", new_rating );

//     let product_review_test = ProductReview{ ProductId: "Dummy".to_owned(), Rating: 1.0, RatingCount: 1.0, SellerDID: "DID".to_owned() };
//     let seller_review_test = SellerReview{ DID: "DID".to_owned(), SellerRating: 2.0, ProductCount: 1.0 };
//     let cbor_product_review:Vec<u8> = serde_cbor::to_vec(&product_review).expect("Failed to serialize to CBOR");
//     let cbor_seller_review:Vec<u8> = serde_cbor::to_vec(&seller_review).expect("Failed to serialize to CBOR");
//     let latest_product_len = cbor_product_review.len();
//     let latest_seller_len = cbor_seller_review.len();
//   //  current_seller_rating = ()/(total_seller_product_count);
//     // append two vectors
//    // let combined_vec = [cbor_product_review, cbor_seller_review].concat();

//     // dump output data
//     unsafe {
//         dump_output(cbor_product_review.as_ptr() , latest_product_len,latest_seller_len);

//     }
// }


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