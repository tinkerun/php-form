import './wasm_exec.js'
import load from './php-form.wasm'

interface Input {
  label: string
  name: string
  value: string
  type: string
}

interface Form {
  (code: string, prefix?: string): Form

  parseCode: (code?: string) => Promise<Input[]>

  stringifyCode: (inputs: string) => Promise<string>
}

declare var PHPFormFunc: Form
declare var Go: any

let phpForm: PHPForm | undefined

export class PHPForm {
  private readonly form: Form

  constructor (form: Form) {
    this.form = form
  }

  async parseCode (code?: string): Promise<Input[]> {
    return await this.form.parseCode(code)
  }

  async stringifyCode (inputs: Input[]): Promise<string> {
    return await this.form.stringifyCode(JSON.stringify(inputs))
  }

  static async instance (code: string, prefix?: string): Promise<PHPForm> {
    if (phpForm != null) {
      return phpForm
    }

    const go = new Go()
    const instance = await load(go.importObject)
    go.run(instance)

    phpForm = new PHPForm(PHPFormFunc(code, prefix))
    return phpForm
  }
}
