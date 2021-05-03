import '../build/wasm_exec.js'

let form

const go = new Go()

WebAssembly.instantiateStreaming(fetch("../build/php-form.wasm"), go.importObject).then((result) => {
    go.run(result.instance)

    form = PHPForm(`
        $next_tinker = [
            'label' => 'TINKER',
            'value' => 'default tinker'
        ];

        $_other = [
            'label' => 'Other',
            'value' => 'default other'
        ];

        echo $next_tinker['value']
    `, 'next')
})

document.querySelector('#btn-parse').addEventListener('click', function () {
    form.parse().then(res => console.log(res))
})

document.querySelector('#btn-stringify').addEventListener('click', function () {
    form.stringify(JSON.stringify([{
        name: '$next_tinker',
        value: 'hello'
    }])).then(res => console.log(res))
})


