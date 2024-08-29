use std::collections::HashMap;
use std::ffi::{CStr, CString};
use std::os::raw::c_char;


#[repr(C)]
#[allow(non_camel_case_types)]
#[derive(Copy, Clone)]
pub struct hmap {
  count: usize,
  entries: *mut hmap_entry,
}

#[repr(C)]
#[allow(non_camel_case_types)]
#[derive(Copy, Clone)]
pub struct hmap_entry {
  key: *mut c_char,
  value: *mut c_char,
}

pub fn hmap_new_from_hashmap(map: &HashMap<String, String>) -> hmap {
  let mut acc = Vec::with_capacity(map.len());

  for (key, value) in map {
    let key = CString::new(key.as_str()).unwrap();
    let value = CString::new(value.as_str()).unwrap();
    acc.push(hmap_entry {
      key: key.into_raw(),
      value: value.into_raw(),
    });
  }

  let count = acc.len();
  let entries = acc.as_mut_ptr();
  std::mem::forget(acc);

  hmap { count, entries }
}

pub unsafe fn hmap_to_hashmap(this: *mut hmap) -> HashMap<String, String> {
  let mut map = HashMap::new();

  if !this.is_null() {
    let hmap = &*this;
    let entries = std::slice::from_raw_parts(hmap.entries, hmap.count);

    for entry in entries {
      let key = CStr::from_ptr(entry.key).to_string_lossy().into_owned();
      let value = CStr::from_ptr(entry.value).to_string_lossy().into_owned();
      map.insert(key, value);
    }
  }

  map
}

#[no_mangle]
pub extern "C" fn hmap_free(this: *mut hmap) {
  if !this.is_null() {
    unsafe {
      let hmap = Box::from_raw(this);
      let entries = std::slice::from_raw_parts_mut(hmap.entries, hmap.count);

      for entry in entries.iter() {
        if !entry.key.is_null() {
          let _ = CString::from_raw(entry.key);
          // dropped automatically
        }

        if !entry.value.is_null() {
          let _ = CString::from_raw(entry.value);
          // dropped automatically
        }
      }

      let _ = Box::from_raw(entries.as_mut_ptr());
      // dropped automatically
    }
  }
}
