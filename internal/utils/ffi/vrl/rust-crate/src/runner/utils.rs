use std::collections::{BTreeMap, HashMap};

use vrl::prelude::*;


pub fn flatten_obj(obj: &BTreeMap<KeyString, Value>) -> HashMap<String, String> {
  let mut result = HashMap::new();

  for (key, value) in obj.iter() {
    match value {
      Value::Bytes(_) => {
        if let Some(s) = value.as_str() {
          result.insert(key.to_string(), s.to_string());
        }
        else {
          result.insert(key.to_string(), "#ERROR#".to_string());
        }
      },
      Value::Regex(r) => {
        result.insert(key.to_string(), r.as_str().to_string());
      },
      Value::Timestamp(ts) => {
        result.insert(key.to_string(), ts.to_string());
      },
      Value::Integer(i) => {
        result.insert(key.to_string(), i.to_string());
      },
      Value::Float(f) => {
        result.insert(key.to_string(), f.to_string());
      },
      Value::Boolean(b) => {
        result.insert(key.to_string(), b.to_string());
      },
      Value::Object(obj) => {
        let inner = flatten_obj(obj);

        for (inner_key, inner_value) in inner.iter() {
          result.insert(format!("{}.{}", key, inner_key), inner_value.to_string());
        }
      },
      Value::Array(arr) => {
        let inner = flatten_array(arr);

        for (inner_key, inner_value) in inner.iter() {
          result.insert(format!("{}.{}", key, inner_key), inner_value.to_string());
        }
      },
      Value::Null => {
        result.insert(key.to_string(), "".to_string());
      },
    };
  }

  result
}

fn flatten_array(array: &Vec<Value>) -> HashMap<String, String> {
  let mut result = HashMap::new();

  for (key, value) in array.iter().enumerate() {
    match value {
      Value::Bytes(_) => {
        if let Some(s) = value.as_str() {
          result.insert(key.to_string(), s.to_string());
        }
        else {
          result.insert(key.to_string(), "#ERROR#".to_string());
        }
      },
      Value::Regex(r) => {
        result.insert(key.to_string(), r.as_str().to_string());
      },
      Value::Timestamp(ts) => {
        result.insert(key.to_string(), ts.to_string());
      },
      Value::Integer(i) => {
        result.insert(key.to_string(), i.to_string());
      },
      Value::Float(f) => {
        result.insert(key.to_string(), f.to_string());
      },
      Value::Boolean(b) => {
        result.insert(key.to_string(), b.to_string());
      },
      Value::Object(obj) => {
        let inner = flatten_obj(obj);

        for (inner_key, inner_value) in inner.iter() {
          result.insert(format!("{}.{}", key, inner_key), inner_value.to_string());
        }
      },
      Value::Array(arr) => {
        let inner = flatten_array(arr);

        for (inner_key, inner_value) in inner.iter() {
          result.insert(format!("{}.{}", key, inner_key), inner_value.to_string());
        }
      },
      Value::Null => {
        result.insert(key.to_string(), "".to_string());
      },
    };
  }

  result
}
