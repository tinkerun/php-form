# php-form

a lib that can modify the php code via javascript

- it's built as [Webassembly](https://webassembly.org/)
- it's written in [Go](https://golang.org) and uses [z7zmey/php-parser](https://github.com/z7zmey/php-parser)

## Install

```
yarn add php-form
```

## API

The full API for php-form is contained within the [TypeScript declaration file](./typings/php-form.d.ts) 

## Example Usage

```js

import {instance} from 'php-form'

(async () => {

  const form = await instance()

  let code = `<?php
  $form_email = [
    'label' => 'Email',
    'value' => 'user1@example.com',
  ];
  
  $form_name = 'billy'
  `
  
  const fields = await form.parse(code)
  // [{name: '$form_email', label: 'Email', value: 'user1@example.com'}, {name: '$form_name', value: 'billy'}]
  
  fields[0].value = 'user2@example.com'
  fields[1].value = 'magic'
  
  code = await form.stringify(fields)
  // <?php $form_email = [
  //   'label' => 'Email',
  //   'value' => 'user2@example.com',
  // ];
  // 
  // $form_name = 'magic'
})


```

## License
[MIT](./LICENSE)