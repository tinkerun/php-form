import './lib/wasm_exec.js'
// @ts-expect-error
import load from './lib/php-form.wasm'

interface PHPFormField {
  name: string
  value: string
}

interface IPHPFormFunc {
  (prefix?: string): IPHPFormFunc
  parse: (code: string) => Promise<PHPFormField[]>
  stringify: (fieldsString: string) => Promise<string>
}

declare var PHPFormFunc: IPHPFormFunc
declare var Go: any

export class PHPForm {
  private readonly form: IPHPFormFunc

  constructor (form: IPHPFormFunc) {
    this.form = form
  }

  async parse (code: string): Promise<PHPFormField[]> {
    return await this.form.parse(code)
  }

  async stringify (fields: PHPFormField[]): Promise<string> {
    return await this.form.stringify(JSON.stringify(fields))
  }
}

let loaded = false

export async function instance (prefix?: string): Promise<PHPForm> {
  if (!loaded) {
    const go = new Go()
    const instance = await load(go.importObject)
    go.run(instance)
    loaded = true
  }

  return new PHPForm(PHPFormFunc(prefix))
}
