import * as request from '@/lib/api/request'
import { InvalidArgumentError } from '@/lib/api/errors'

type TestTransformerResult =
  | { success: true, record: Record<string, string> }
  | { success: false, error: string }

export const testTransformer = async (
  code: string,
  record: Record<string, string>
): Promise<TestTransformerResult> => {
  type TestTransformerRequest = {
    code: string
    record: Record<string, string>
  }

  type TestTransformerResponse = {
    success: boolean
    record: Record<string, string>
  }

  try {
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
    return { success: true, record: body.record }
  }
  catch (err) {
    if (err instanceof InvalidArgumentError && err.appcode === 422) {
      return { success: false, error: err.message }
    }
    else {
      throw err
    }
  }
}
