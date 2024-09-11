use std::collections::{BTreeMap, HashMap};

use anyhow::{Context, Result};
use vrl::prelude::*;


pub fn process_record(
  record: HashMap<String, String>,
  script: String,
) -> Result<HashMap<String, String>> {
  let mut vrl_record = ObjectMap::new();
  for (key, value) in record.iter() {
    vrl_record.insert(key.clone().into(), Value::from(value.clone()));
  }

  let fns = vrl::stdlib::all();
  let compiled = match vrl::compiler::compile(&script, &fns) {
    Ok(compiled) => compiled,
    Err(diagnostics) => {
      let messages: Vec<String> = diagnostics.iter()
        .map(|d| d.message.clone())
        .collect();

      return Err(anyhow::anyhow!("Failed to compile VRL script"))
        .context(messages.join("\n"))
    }
  };

  let mut target = vrl::compiler::TargetValue {
    value: Value::Object(vrl_record),
    metadata: Value::Object(ObjectMap::new()),
    secrets: vrl::value::Secrets::default(),
  };

  let mut state = state::RuntimeState::default();
  let timezone = TimeZone::default();
  let mut ctx = vrl::compiler::Context::new(&mut target, &mut state, &timezone);

  let vrl_result = compiled.program.resolve(&mut ctx)
    .context("Failed to execute VRL script")?;

  let result = if let Value::Object(obj) = vrl_result {
    flatten_obj(&obj)
  } else {
    HashMap::new()
  };

  Ok(result)
}

fn flatten_obj(obj: &BTreeMap<KeyString, Value>) -> HashMap<String, String> {
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

#[cfg(test)]
mod tests {
  use super::*;

  #[test]
  fn test_process_record() {
    let input = HashMap::new();
    let script = String::from(r#"
      .foo = "x"
      .bar.baz = [1, 2, 3, "a"]
      .
    "#);

    let output = process_record(input, script).unwrap();

    assert_eq!(output.get("foo"), Some(&"x".to_string()));
    assert_eq!(output.get("bar.baz.0"), Some(&"1".to_string()));
    assert_eq!(output.get("bar.baz.1"), Some(&"2".to_string()));
    assert_eq!(output.get("bar.baz.2"), Some(&"3".to_string()));
    assert_eq!(output.get("bar.baz.3"), Some(&"a".to_string()));
  }
}
