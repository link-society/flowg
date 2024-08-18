use std::ffi::{CString, CStr};
use std::os::raw::c_char;

mod compiler;

#[repr(C)]
#[allow(non_camel_case_types)]
pub struct filterdsl_result {
  success: bool,
  data: *mut c_char,
}

#[no_mangle]
pub extern "C" fn filterdsl_compile(input: *const c_char) -> filterdsl_result {
  if input.is_null() {
    let err_msg = CString::new("input is null").unwrap();
    return filterdsl_result {
      success: false,
      data: err_msg.into_raw(),
    };
  }

  let c_str = unsafe { CStr::from_ptr(input) };

  let input_str = c_str.to_str().unwrap_or_default().to_string();
  match compiler::compile(input_str) {
    Ok(json_data) => {
      let json_data = CString::new(json_data).unwrap();

      filterdsl_result {
        success: true,
        data: json_data.into_raw(),
      }
    },
    Err(err) => {
      let err_msg = CString::new(err.to_string()).unwrap();

      filterdsl_result {
        success: false,
        data: err_msg.into_raw(),
      }
    }
  }
}

#[no_mangle]
pub extern "C" fn filterdsl_result_free(result: filterdsl_result) {
  if result.data.is_null() {
    return;
  }

  unsafe {
    let _ = CString::from_raw(result.data);
  }
}
