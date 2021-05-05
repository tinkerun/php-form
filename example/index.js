import {PHPForm} from '../lib/index.js';

(async () => {
    const form = await PHPForm.instance(`
    $next_tinker = [
        'label' => 'TINKER',
        'value' => 'default tinker'
    ];
    
    $_other = [
        'label' => 'Other',
        'value' => 'default other'
    ];
    
    echo $next_tinker['value']
    `, 'next');

    document.querySelector('#btn-parse').addEventListener('click', async () => {
        const res = await form.parseCode()
        console.log(res)
    })
    
    document.querySelector('#btn-stringify').addEventListener('click', function () {
        form.stringifyCode([{
            name: '$next_tinker',
            value: 'hello'
        }]).then(res => console.log(res))
    })
    
})()
