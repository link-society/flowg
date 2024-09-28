use std::collections::HashMap;
use std::ffi::CString;
use std::os::raw::c_char;

use super::{hmap, hmap_new_from_hashmap, hmap_free};

#[repr(C)]
#[allow(non_camel_case_types)]
pub struct vrl_result {
  tag: vrl_result_tag,
  data: vrl_result_data,
}

#[repr(C)]
#[allow(non_camel_case_types)]
pub enum vrl_result_tag {
  vrl_result_ok,
  vrl_result_err,
}

#[repr(C)]
#[allow(non_camel_case_types)]
pub struct vrl_result_data {
  ok_data: *mut hmap,
  err_data: *mut c_char,
}


pub fn vrl_result_new_ok(map: HashMap<String, String>) -> *mut vrl_result {
  let result = vrl_result {
    tag: vrl_result_tag::vrl_result_ok,
    data: vrl_result_data {
      ok_data: Box::into_raw(Box::new(hmap_new_from_hashmap(&map))),
      err_data: std::ptr::null_mut(),
    },
  };
  Box::into_raw(Box::new(result))
}

pub fn vrl_result_new_err(err: &str) -> *mut vrl_result {
  let result = vrl_result {
    tag: vrl_result_tag::vrl_result_err,
    data: vrl_result_data {
      ok_data: std::ptr::null_mut(),
      err_data: CString::new(err).unwrap().into_raw(),
    },
  };
  Box::into_raw(Box::new(result))
}

#[no_mangle]
pub extern "C" fn vrl_result_free(this: *mut vrl_result) {
  if !this.is_null() {
    unsafe {
      let result = Box::from_raw(this);

      match result.tag {
        vrl_result_tag::vrl_result_ok => {
          hmap_free(result.data.ok_data);
        },
        vrl_result_tag::vrl_result_err => {
          let _ = CString::from_raw(result.data.err_data);
          // dropped automatically
        },
      }
    }
  }
}
