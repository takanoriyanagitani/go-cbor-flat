use std::sync::RwLock;

static INPUT: RwLock<Vec<u8>> = RwLock::new(vec![]);
static OUTPUT: RwLock<Vec<u8>> = RwLock::new(vec![]);

pub fn _initialize_input_buffer(size: u32, init: u8) -> Result<u32, &'static str> {
    let mut guard = INPUT.try_write().map_err(|_| "unable to write lock")?;
    let mv: &mut Vec<_> = &mut guard;
    mv.resize(size as usize, init);
    Ok(mv.capacity() as u32)
}

#[allow(unsafe_code)]
#[no_mangle]
pub extern "C" fn initialize_input_buffer(size: u32, init: u8) -> i32 {
    _initialize_input_buffer(size, init)
        .ok()
        .and_then(|u| u.try_into().ok())
        .unwrap_or(-1)
}

pub fn _get_input_offset() -> Result<*mut u8, &'static str> {
    let mut guard = INPUT.try_write().map_err(|_| "unable to write lock")?;
    let mv: &mut Vec<_> = &mut guard;
    Ok(mv.as_mut_ptr())
}

#[allow(unsafe_code)]
#[no_mangle]
pub fn get_input_offset() -> *mut u8 {
    _get_input_offset().ok().unwrap_or_else(std::ptr::null_mut)
}

pub fn copy(i: &[u8], o: &mut Vec<u8>) -> usize {
    o.resize(i.len(), 0);
    o.clear();
    o.extend(i);
    o.len()
}

pub fn _convert() -> Result<u32, &'static str> {
    let guard = INPUT.try_read().map_err(|_| "unable to read lock")?;
    let i: &[u8] = &guard;

    let mut guard = OUTPUT.try_write().map_err(|_| "unable to write lock")?;
    let o: &mut Vec<_> = &mut guard;

    let sz: usize = copy(i, o);
    Ok(sz as u32)
}

#[allow(unsafe_code)]
#[no_mangle]
pub fn convert() -> i32 {
    _convert()
        .ok()
        .and_then(|u| u.try_into().ok())
        .unwrap_or(-1)
}

pub fn _get_output_offset() -> Result<*const u8, &'static str> {
    let guard = INPUT.try_read().map_err(|_| "unable to read lock")?;
    let v: &Vec<_> = &guard;
    Ok(v.as_ptr())
}

#[allow(unsafe_code)]
#[no_mangle]
pub fn get_output_offset() -> *const u8 {
    _get_output_offset().ok().unwrap_or_else(std::ptr::null)
}
