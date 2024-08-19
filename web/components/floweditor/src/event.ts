export class AddNodeEvent extends CustomEvent<{ type: string }> {
  constructor(type: string) {
    super('add-node', { detail: { type } })
  }
}
