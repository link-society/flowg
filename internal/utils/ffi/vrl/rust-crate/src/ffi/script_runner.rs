use std::ffi::*;

use crate::{
  ffi::hmap::*,
  runner::ScriptRunner,
};

#[no_mangle]
pub extern "C" fn script_runner_new(
  c_source: *const c_char,
  c_err: *mut *mut c_char,
) -> *mut c_void {
  let source = unsafe {
    assert!(!c_source.is_null());
    CStr::from_ptr(c_source).to_string_lossy().into_owned()
  };

  match ScriptRunner::new(&source) {
    Ok(runner) => {
      let boxed = Box::new(runner);
      Box::into_raw(boxed) as *mut c_void
    }
    Err(e) => {
      if !c_err.is_null() {
        let err_msg = CString::new(format!("{:#}", e)).unwrap();
        unsafe {
          *c_err = err_msg.into_raw();
        }
      }
      std::ptr::null_mut()
    }
  }
}

#[no_mangle]
pub extern "C" fn script_runner_free(this: *mut c_void) {
  if !this.is_null() {
    unsafe {
      let _ = Box::from_raw(this as *mut ScriptRunner);
      // dropped automatically
    }
  }
}

#[no_mangle]
pub extern "C" fn script_runner_eval(
  this: *mut c_void,
  c_input: *mut hmap,
  c_err: *mut *mut c_char,
) -> *mut hmap {
  let runner = unsafe {
    assert!(!this.is_null());
    &*(this as *mut ScriptRunner)
  };
  let input = hmap_to_hashmap(c_input);

  match runner.process_record(input) {
    Ok(output) => {
      let boxed = Box::new(hmap_new_from_hashmap(&output));
      Box::into_raw(boxed)
    },
    Err(e) => {
      if !c_err.is_null() {
        let err_msg = CString::new(format!("{:#}", e)).unwrap();
        unsafe {
          *c_err = err_msg.into_raw();
        }
      }
      std::ptr::null_mut()
    }
  }
}
