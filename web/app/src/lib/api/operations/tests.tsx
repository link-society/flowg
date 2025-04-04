import * as request from '@/lib/api/request'

export const testTransformer = async (
  code: string,
  record: Record<string, string>
): Promise<Record<string, string>> => {
  type TestTransformerRequest = {
    code: string
    record: Record<string, string>
  }

  type TestTransformerResponse = {
    success: boolean
    record: Record<string, string>
  }

  const { body } = await request.POST<
    TestTransformerRequest,
    TestTransformerResponse
  >({
    path: '/api/v1/test/transformer',
    body: {
      code,
      record,
    },
  })
  return body.record
}
