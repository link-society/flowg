use ariadne::{Config, Label, Report, ReportKind, Source};
use vrl::{
  diagnostic::{DiagnosticList, Severity},
  compiler::ExpressionError,
};

const FILENAME: &str = "<script>";

pub fn render_diagnostics(
  diagnostics: DiagnosticList,
  script: &str,
) -> String {
  let mut out = Vec::new();
  let cache = (FILENAME, Source::from(script));

  for diagnostic in diagnostics.iter() {
    let kind = match diagnostic.severity {
      Severity::Bug | Severity::Error => ReportKind::Error,
      Severity::Warning => ReportKind::Warning,
      Severity::Note => ReportKind::Advice,
    };

    let span = diagnostic.labels
      .iter()
      .find(|label| label.primary)
      .or_else(|| diagnostic.labels.first())
      .map(|label| (FILENAME, label.span.start()..label.span.end()))
      .unwrap_or((FILENAME, 0..0));

    let mut builder = Report::build(kind, span)
      .with_config(
        Config::default()
          .with_color(false)
          .with_compact(true)
      )
      .with_code(diagnostic.code)
      .with_message(&diagnostic.message)
      .with_labels(diagnostic.labels.iter().map(|label| {
        Label::new((FILENAME, label.span.start()..label.span.end()))
          .with_message(&label.message)
      }));

    builder.with_notes(diagnostic.notes.iter().map(|note| note.to_string()));

    builder.finish().write(cache.clone(), &mut out).expect("Failed to write diagnostic report");

  }

  String::from_utf8_lossy(&out).to_string()
}

pub fn render_error(
  error: ExpressionError,
  script: &str,
) -> String {
  render_diagnostics(DiagnosticList::from(error), script)
}
