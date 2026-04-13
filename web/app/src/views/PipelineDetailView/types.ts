import PipelineModel from '@/lib/models/PipelineModel'

export type LoaderData = {
  pipelines: string[]
  currentPipeline: {
    name: string
    flow: PipelineModel
  }
}
