use std::collections::HashMap;

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

  let mut result = HashMap::new();

  if let Value::Object(obj) = vrl_result {
    for (key, value) in obj.iter() {
      if let Some(s) = value.as_str() {
        result.insert(key.to_string(), s.to_string());
      }
    }
  }

  Ok(result)
}
