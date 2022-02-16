package zasm

import (
	"bytes"
	"testing"

	. "github.com/ofabricio/calm"
	"github.com/stretchr/testify/assert"
)

func TestPrintDirective(t *testing.T) {

	var buf bytes.Buffer

	p := &Parser{c: New(string(src))}
	ast, _ := p.Parse()
	Walk(&printDirectiveVisitor{buf: &buf}, ast)

	assert.Equal(t, exp, buf.String())
}

const src = `
@print {
}

@print {
	@print {
		@print {
		}
	}
}

@print {
	mov eax eax
}

@print {
	mov eax eax
	@print {
		mov eax eax
		@print {
			mov eax eax
		}
		mov eax eax
	}
	mov eax eax
}

@print {
	mov eax eax
	@print {
		mov eax eax
		@print {
			mov eax eax
		}
	}
	mov eax eax
}

@print {
	mov eax eax
	@print {
		@print {
			mov eax eax
		}
		mov eax eax
	}
	mov eax eax
}

@print {
	mov eax eax
	@print {
		mov eax eax
		@print {
		}
		mov eax eax
	}
	mov eax eax
}

`

const exp = `

    [2 bytes]
    
    89 C0    mov eax eax

    [10 bytes]
    
    89 C0    mov eax eax
    
    |   [6 bytes]
    |   
    |   89 C0    mov eax eax
    |   
    |   |   [2 bytes]
    |   |   
    |   |   89 C0    mov eax eax
    |   
    |   89 C0    mov eax eax
    
    89 C0    mov eax eax

    [8 bytes]
    
    89 C0    mov eax eax
    
    |   [4 bytes]
    |   
    |   89 C0    mov eax eax
    |   
    |   |   [2 bytes]
    |   |   
    |   |   89 C0    mov eax eax
    
    89 C0    mov eax eax

    [8 bytes]
    
    89 C0    mov eax eax
    
    |   [4 bytes]
    |   
    |   |   [2 bytes]
    |   |   
    |   |   89 C0    mov eax eax
    |   
    |   89 C0    mov eax eax
    
    89 C0    mov eax eax

    [8 bytes]
    
    89 C0    mov eax eax
    
    |   [4 bytes]
    |   
    |   89 C0    mov eax eax
    |   89 C0    mov eax eax
    
    89 C0    mov eax eax

`
