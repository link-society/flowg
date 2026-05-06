mod report;

use std::sync::Mutex;
use vrl::prelude::*;
use anyhow::{Context, Result};

pub struct ScriptRunner {
  inner: Mutex<ScriptRunnerInner>,

  pub outbuf: Vec<u8>,
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
      outbuf: Vec::with_capacity(1024),
    })
  }

  pub fn eval(&self, input: Value) -> Result<Value> {
    let guard = self.inner.lock().unwrap();

    let mut target = vrl::compiler::TargetValue {
      value: input,
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

    Ok(target.value)
  }
}

#[cfg(test)]
mod tests {
  use super::*;

  #[test]
  fn test_input() {
    let input = Value::Object(ObjectMap::new());
    let script = r#"
      .foo = "x"
      .bar.baz = [1, 2, 3, "a"]
      .
    "#;

    let runner = ScriptRunner::new(script).unwrap();
    let output = runner.eval(input).unwrap();

    assert!(output.is_object());

    let obj = output.as_object().unwrap();
    assert_eq!(obj.get("foo"), Some(&Value::Bytes(Bytes::from("x"))));

    assert!(obj.get("bar").is_some());
    assert!(obj.get("bar").unwrap().is_object());

    let bar = obj.get("bar").unwrap().as_object().unwrap();
    assert!(bar.get("baz").is_some());
    assert!(bar.get("baz").unwrap().is_array());

    let baz = bar.get("baz").unwrap().as_array().unwrap();
    assert_eq!(baz.len(), 4);
    assert_eq!(baz[0], Value::Integer(1));
    assert_eq!(baz[1], Value::Integer(2));
    assert_eq!(baz[2], Value::Integer(3));
    assert_eq!(baz[3], Value::Bytes(Bytes::from("a")));
  }
}
