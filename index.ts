import './lib/wasm_exec.js'
// @ts-expect-error
import load from './lib/php-form.wasm'

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

declare global {
  interface Window {
    PHPFormFunc: Form
    Go: any
  }
}

export class PHPForm {
  private static phpForm: PHPForm | undefined

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
    if (PHPForm.phpForm != null) {
      return PHPForm.phpForm
    }

    const go = new window.Go()
    const instance = await load(go.importObject)
    go.run(instance)

    PHPForm.phpForm = new PHPForm(window.PHPFormFunc(code, prefix))
    return PHPForm.phpForm
  }
}
