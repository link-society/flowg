use logos::{Lexer, Logos, Span};

#[derive(Logos, Debug, PartialEq)]
#[logos(skip r"[ \t\n\f]+")]
pub enum Token {
  #[regex(r"[a-zA-Z_][a-zA-Z0-9_]*", |lex| lex.slice().to_string())]
  #[regex(r"'([^'\\]|\\.)*'", unescape)]
  Ident(String),

  #[regex(r#""([^"\\]|\\.)*""#, unescape)]
  String(String),

  #[token("(")]
  LParen,

  #[token(")")]
  RParen,

  #[token("[")]
  LBracket,

  #[token("]")]
  RBracket,

  #[token(",")]
  Comma,

  #[token("=")]
  Eq,

  #[token("!=")]
  Ne,

  #[token("NOT")]
  #[token("not")]
  OpNot,

  #[token("AND")]
  #[token("and")]
  OpAnd,

  #[token("OR")]
  #[token("or")]
  OpOr,

  #[token("IN")]
  #[token("in")]
  OpIn,
}

fn unescape<'a>(lex: &mut Lexer<'a, Token>) -> String {
  snailquote::unescape(lex.slice()).unwrap()
}

#[derive(Debug, Clone, Copy)]
pub struct TokenLocation(pub usize, pub usize);

impl std::fmt::Display for TokenLocation {
  fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
    write!(f, "{}..{}", self.0, self.1)
  }
}

pub struct TokenStream {
  pub size: usize,
  pub tokens: Vec<(Token, Span)>,
}

impl peg::Parse for TokenStream {
  type PositionRepr = TokenLocation;

  fn start(&self) -> usize {
    0
  }

  fn is_eof(&self, pos: usize) -> bool {
    pos >= self.tokens.len()
  }

  fn position_repr(&self, pos: usize) -> Self::PositionRepr {
    match self.tokens.get(pos) {
      Some((_token, span)) => TokenLocation(span.start, span.end),
      None => TokenLocation(self.size, self.size),
    }
  }
}

impl<'src> peg::ParseElem<'src> for TokenStream {
  type Element = &'src Token;

  fn parse_elem(&'src self, pos: usize) -> peg::RuleResult<Self::Element> {
    match self.tokens.get(pos) {
      Some((token, _)) => peg::RuleResult::Matched(pos + 1, token),
      None => peg::RuleResult::Failed,
    }
  }
}

impl peg::ParseLiteral for TokenStream {
  fn parse_string_literal(&self, pos: usize, literal: &str) -> peg::RuleResult<()> {
    match (literal, self.tokens.get(pos)) {
      _ => peg::RuleResult::Failed,
    }
  }
}

impl<'src> peg::ParseSlice<'src> for TokenStream {
  type Slice = Vec<&'src Token>;

  fn parse_slice(&'src self, begin: usize, end: usize) -> Self::Slice {
    self.tokens[begin..end]
      .into_iter()
      .map(|(token, _)| token)
      .collect()
  }
}
