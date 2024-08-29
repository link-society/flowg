mod hmap_struct;
mod result_struct;

pub use self::hmap_struct::{
  hmap,
  hmap_entry,
  hmap_new_from_hashmap,
  hmap_to_hashmap,
  hmap_free,
};
pub use self::result_struct::{
  vrl_result,
  vrl_result_tag,
  vrl_result_data,
  vrl_result_new_ok,
  vrl_result_new_err,
  vrl_result_free,
};
