use std::ffi::{CStr, c_char};

mod runner;
mod ffi;

pub use ffi::{
  hmap,
  hmap_entry,
  hmap_new_from_hashmap,
  hmap_to_hashmap,
  hmap_free,
  vrl_result,
  vrl_result_tag,
  vrl_result_data,
  vrl_result_new_ok,
  vrl_result_new_err,
  vrl_result_free,
};

#[no_mangle]
pub extern "C" fn process_record(
  input: *mut hmap,
  script: *const c_char,
) -> *mut vrl_result {
  if input.is_null() || script.is_null() {
    return vrl_result_new_err("Invalid input: null pointer received.");
  }

  let script = unsafe { CStr::from_ptr(script).to_str().unwrap_or_default().to_string() };

  let record = unsafe { hmap_to_hashmap(input) };

  match runner::process_record(record, script) {
    Ok(result) => vrl_result_new_ok(result),
    Err(e) => vrl_result_new_err(format!("{}", e).as_str()),
  }
}
