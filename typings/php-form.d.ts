declare module 'php-form' {
    export interface PHPForm {
        /**
         * @param 需要生成表单的代码
         * @param 自定义变量前缀，默认使用 `form_`
         */
        (code: string, prefix?: string): PHPForm

        /**
         * 解析 php 代码，得到表单结构
         * @param code 需要解析的代码，默认使用初始化时候的 code
         */
        parseCode: (code?: string) => Promise<PHPFormInput[]>

        /**
         * 根据表单的输入数据，返回填充值之后的代码
         * @param 表单的输入数据
         */
        stringifyCode: (inputs: string | PHPFormInput[] | Uint8Array) => Promise<string>
    }

    export interface PHPFormInput {
        /**
         * 表单字段的名称，可以理解为 `<input/>` 的 `name` 属性
         */
        name: string

        /**
         * 表单字段的值，可以理解为 `<input/>` 的 `value` 属性
         */
        value: string

        /**
         * 表单字段的标签，可以理解为 `<label/>`
         */
        label?: string

        /**
         * 表单字段的类型，可以理解为 `<input/>` 的 `type` 属性
         */
        type?: string
    }
}
