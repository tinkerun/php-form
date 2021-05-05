declare module 'php-form' {
    /**
    * 获取并且初始化 PHPForm 对象
    * 
    * @param code 需要生成表单的代码
    * @param prefix 自定义变量前缀，默认使用 `form_`
    */
    export function instance (code?: string, prefix?: string): Promise<PHPForm>

    export interface PHPForm {
        /**
         * 解析 php 代码，得到表单结构
         * @param code 需要解析的代码，默认使用初始化时候的 code
         */
        parse: (code?: string) => Promise<PHPFormField[]>

        /**
         * 根据表单的输入数据，返回填充值之后的代码
         * @param 表单的输入数据
         */
        stringify: (fields: PHPFormField[]) => Promise<string>
    }

    export interface PHPFormField {
        /**
         * 表单字段的名称，可以理解为 `<input/>` 的 `name` 属性
         */
        name: string

        /**
         * 表单字段的值，可以理解为 `<input/>` 的 `value` 属性
         */
        value: string
    }
}
