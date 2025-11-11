mod utils;
mod report;

use std::{collections::HashMap, sync::Mutex};
use vrl::prelude::*;
use anyhow::{Context, Result};

pub struct ScriptRunner {
  inner: Mutex<ScriptRunnerInner>,
}

struct ScriptRunnerInner {
  program: vrl::compiler::Program,
  source: String,
}

impl ScriptRunner {
  pub fn new(source: &str) -> Result<Self> {
    let fns = vrl::stdlib::all();
    let program = match vrl::compiler::compile(source, &fns) {
      Ok(compiled) => compiled.program,
      Err(diagnostics) => {
        return Err(anyhow::anyhow!("Failed to compile VRL script"))
          .context(report::render_diagnostics(diagnostics, source))
      }
    };

    Ok(ScriptRunner {
      inner: Mutex::new(ScriptRunnerInner {
        program,
        source: source.to_string(),
      }),
    })
  }

  pub fn process_record(
    &self,
    record: HashMap<String, String>,
  ) -> Result<HashMap<String, String>> {
    let guard = self.inner.lock().unwrap();

    let mut obj = ObjectMap::new();
    for (key, value) in record.iter() {
      obj.insert(key.clone().into(), Value::from(value.clone()));
    }

    let mut target = vrl::compiler::TargetValue {
      value: Value::Object(obj),
      metadata: Value::Object(ObjectMap::new()),
      secrets: vrl::value::Secrets::default(),
    };

    let mut state = state::RuntimeState::default();
    let timezone = TimeZone::default();
    let mut ctx = vrl::compiler::Context::new(&mut target, &mut state, &timezone);

    if let Err(expression_error) = guard.program.resolve(&mut ctx) {
    return Err(anyhow::anyhow!("Failed to execute VRL script"))
      .context(report::render_error(expression_error, &guard.source));
    }

    let result = if let Value::Object(obj) = target.value {
      utils::flatten_obj(&obj)
    } else {
      HashMap::new()
    };

    Ok(result)
  }
}

#[cfg(test)]
mod tests {
  use super::*;

  #[test]
  fn test_process_record() {
    let input = HashMap::new();
    let script = r#"
      .foo = "x"
      .bar.baz = [1, 2, 3, "a"]
      .
    "#;

    let runner = ScriptRunner::new(script).unwrap();
    let output = runner.process_record(input).unwrap();

    assert_eq!(output.get("foo"), Some(&"x".to_string()));
    assert_eq!(output.get("bar.baz.0"), Some(&"1".to_string()));
    assert_eq!(output.get("bar.baz.1"), Some(&"2".to_string()));
    assert_eq!(output.get("bar.baz.2"), Some(&"3".to_string()));
    assert_eq!(output.get("bar.baz.3"), Some(&"a".to_string()));
  }
}
