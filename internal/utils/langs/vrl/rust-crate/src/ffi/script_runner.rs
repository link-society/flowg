use std::ffi::*;

use crate::{
  runner::ScriptRunner,
};

#[allow(non_camel_case_types)]
pub type script_runner = *mut c_void;

#[repr(C)]
pub struct compilation_result {
  pub is_ok: bool,
  pub ok: compilation_result_ok_variant,
  pub err: compilation_result_err_variant,
}

#[repr(C)]
pub struct compilation_result_ok_variant {
  pub runner: script_runner,
}

#[repr(C)]
pub struct compilation_result_err_variant {
  pub reason: *mut c_char,
}


#[no_mangle]
pub extern "C" fn compile_script(
  c_source: *const c_char,
) -> compilation_result {
  let source = unsafe {
    assert!(!c_source.is_null());
    CStr::from_ptr(c_source).to_string_lossy().into_owned()
  };

  match ScriptRunner::new(&source) {
    Ok(runner) => {
      let boxed = Box::new(runner);
      let handle = Box::into_raw(boxed) as *mut c_void;

      compilation_result {
        is_ok: true,
        ok: compilation_result_ok_variant {
          runner: handle,
        },
        err: compilation_result_err_variant {
          reason: std::ptr::null_mut(),
        },
      }
    }
    Err(e) => {
      compilation_result {
        is_ok: false,
        ok: compilation_result_ok_variant {
          runner: std::ptr::null_mut(),
        },
        err: compilation_result_err_variant {
          reason: CString::new(format!("{:#}", e)).unwrap().into_raw(),
        },
      }
    }
  }
}

#[no_mangle]
pub extern "C" fn drop_script_runner(this: script_runner) {
  if !this.is_null() {
    unsafe {
      let _ = Box::from_raw(this as *mut ScriptRunner);
      // dropped automatically
    }
  }
}
