use std::ffi::*;

use crate::ffi::{
  msgpack_buffer::*,
  script_runner::*,
};
use crate::runner::ScriptRunner;

use vrl::value::Value;


#[repr(C)]
pub struct log_transform_result {
  is_ok: bool,
  ok: log_transform_result_ok_variant,
  err: log_transform_result_err_variant,
}

#[repr(C)]
pub struct log_transform_result_ok_variant {
  data: msgpack_buffer,
}

#[repr(C)]
pub struct log_transform_result_err_variant {
  reason: *mut c_char,
}


#[no_mangle]
pub extern "C" fn transform_log(
  runner: script_runner,
  log_event: msgpack_buffer,
) -> log_transform_result {
  let log_event = unsafe {
    assert!(!log_event.data.is_null());
    assert!(log_event.len > 0);
    std::slice::from_raw_parts(log_event.data, log_event.len)
  };
  let log_event: Value = match rmp_serde::from_slice(log_event) {
    Ok(v) => v,
    Err(e) => {
      return log_transform_result {
        is_ok: false,
        ok: log_transform_result_ok_variant {
          data: msgpack_buffer {
            data: std::ptr::null_mut(),
            len: 0,
          },
        },
        err: log_transform_result_err_variant {
          reason: CString::new(format!("failed to deserialize log event: {:#}", e)).unwrap().into_raw(),
        },
      }
    }
  };

  let runner = unsafe {
    assert!(!runner.is_null());
    &mut *(runner as *mut ScriptRunner)
  };

  let res = runner
    .eval(log_event)
    .map_err(|e| format!("failed to transform log event: {:#}", e))
    .and_then(|v| {
      runner.outbuf.clear();
      rmp_serde::encode::write(&mut runner.outbuf, &v)
        .map_err(|e| format!("failed to serialize transformed log event: {:#}", e))
    });

  match res {
    Ok(_) => log_transform_result {
      is_ok: true,
      ok: log_transform_result_ok_variant {
        data: msgpack_buffer {
          data: runner.outbuf.as_mut_ptr(),
          len: runner.outbuf.len(),
        },
      },
      err: log_transform_result_err_variant {
        reason: std::ptr::null_mut(),
      },
    },
    Err(reason) => log_transform_result {
      is_ok: false,
      ok: log_transform_result_ok_variant {
        data: msgpack_buffer {
          data: std::ptr::null_mut(),
          len: 0,
        },
      },
      err: log_transform_result_err_variant {
        reason: CString::new(reason).unwrap().into_raw(),
      },
    },
  }
}
