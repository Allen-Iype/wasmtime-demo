extern "C" {
    fn load_input(pointer: *mut u8);
    fn dump_output(pointer: *const u8, bill_paid:u32, length: usize);
}

#[no_mangle]
pub extern "C" fn handler(user_id_length: usize , bill_length: usize , balance_length: usize , status_length: usize) {
    // load input data
    let mut input = Vec::with_capacity(user_id_length + bill_length + balance_length + status_length);
    
    unsafe {
        load_input(input.as_mut_ptr());
        input.set_len(user_id_length + bill_length + balance_length + status_length);
    }


    let (user_id, b1_rest) = input.split_at(user_id_length);
    let (bill, balance_status) = b1_rest.split_at(bill_length);
    let (balance, user_status) = balance_status.split_at(balance_length);

    // process app data
    let output = user_id.to_ascii_uppercase();
    println!("Did ouput{:?}",output);

    // let uservote = u32::from_ne_bytes(b1[0..b1_length].try_into().unwrap());
    // let mut redvote = u32::from_ne_bytes(red[0..redcount].try_into().unwrap());
    // let mut bluevote = u32::from_ne_bytes(blue[0..bluecount].try_into().unwrap());
 //   let user_id_value = String::from_utf8(user_id.to_vec()).unwrap();
    let bill_value = u32::from_ne_bytes(bill.try_into().unwrap());
    let balance_value = u32::from_ne_bytes(balance.try_into().unwrap());
 //   let user_status_value = String::from_utf8(user_status.to_vec()).unwrap();
    
    // if uservote == 1 {
    // redvote += 1;
    // }
    // else {
    //     bluevote += 1
    // }
        if bill_value > balance_value {
            let bill_paid = 0;
            unsafe {
                dump_output(output.as_ptr() , bill_paid , output.len());
            }
        }
        else {
            let bill_paid = 1;
            unsafe {
                dump_output(output.as_ptr() , bill_paid , output.len());
            }
        }
    // dump output data
}
