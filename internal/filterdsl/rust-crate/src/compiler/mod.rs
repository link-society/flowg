use anyhow::{anyhow, Result};
use logos::{Logos, Span};

mod lexer;
mod grammar;

use lexer::{Token, TokenLocation, TokenStream};

pub fn compile(input: String) -> Result<String> {
  let tokens = Token::lexer(&input)
    .spanned()
    .map(|(res, span)| match res {
      Ok(token) => Ok((token, span)),
      Err(()) => Err(anyhow!("syntax error: {}", &input[span])),
    })
    .collect::<Result<Vec<(Token, Span)>>>()?;

  let token_stream = TokenStream {
    size: input.len(),
    tokens,
  };

  let json_data = grammar::parser::start(&token_stream)
    .map_err(|err| {
      let TokenLocation(start, end) = err.location;
      anyhow!("unexpected token: {}, expected {}", &input[start..end], err.expected)
    })?;

  Ok(json_data)
}
