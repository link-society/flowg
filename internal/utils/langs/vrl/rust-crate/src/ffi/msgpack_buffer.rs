#[repr(C)]
pub struct msgpack_buffer {
  pub data: *const u8,
  pub len: usize,
}
