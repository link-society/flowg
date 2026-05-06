import { InvalidArgumentError } from '@/lib/api/errors'
import * as request from '@/lib/api/request'

type TestTransformerResult =
  | { success: true; records: Array<Record<string, string>> }
  | { success: false; error: string }

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
    records: Array<Record<string, string>>
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
    return { success: true, records: body.records }
  } catch (err) {
    if (err instanceof InvalidArgumentError && err.appcode === 422) {
      return { success: false, error: err.message }
    } else {
      throw err
    }
  }
}
