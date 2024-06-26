<template>
    <b-row class="mb-5">
        <b-col>
            <div id="terminal" style="height: 1000px"></div>
        </b-col>
    </b-row>
</template>

<script>
    import axios from 'axios';
    import { Terminal } from 'xterm';
    import { FitAddon } from 'xterm-addon-fit';

    export default {
        name: 'Syslog',
        data() {
            return {
                logs: [],
                cur: parseInt(this.$route.query.next !== undefined ? this.$route.query.next : 0),
                next: -1,
            };
        },
        methods: {
            typeWrap(typ) {
                switch(typ) {
                    case 'ERROR':
                        return '\x1B[1;3;31m[错误]\x1B[0m => ';
                    case 'ACTION':
                        return '\x1B[1;3;32m[操作]\x1B[0m => ';
                    case 'SYSTEM':
                        return '\x1B[1;3;33m[系统]\x1B[0m => ';
                }

                return 'UNKNOWN'
            }
        },
        mounted() {
            const term = new Terminal();
            const fitAddon = new FitAddon();
            term.loadAddon(fitAddon);
            term.open(document.getElementById("terminal"));

            let that = this;
            axios.get('/api/syslog/logs/?limit=100&offset=' + this.cur).then(response => {
                this.next = response.data.next;
                this.logs = response.data.logs.reverse().map((line) => {
                    term.writeln(that.typeWrap(line.type) + line.body);
                });
                fitAddon.fit();
            }).catch(error => {
                this.ToastError(error);
            });
        }
    }
</script>

<style>
    .xterm {
        font-feature-settings: "liga" 0;
        position: relative;
        user-select: none;
        -ms-user-select: none;
        -webkit-user-select: none;
    }
    .xterm.focus,
    .xterm:focus {
        outline: none;
    }
    .xterm .xterm-helpers {
        position: absolute;
        top: 0;
        /**
         * The z-index of the helpers must be higher than the canvases in order for
         * IMEs to appear on top.
         */
        z-index: 5;
    }
    .xterm .xterm-helper-textarea {
        /*
         * HACK: to fix IE's blinking cursor
         * Move textarea out of the screen to the far left, so that the cursor is not visible.
         */
        position: absolute;
        opacity: 0;
        left: -9999em;
        top: 0;
        width: 0;
        height: 0;
        z-index: -5;
        /** Prevent wrapping so the IME appears against the textarea at the correct position */
        white-space: nowrap;
        overflow: hidden;
        resize: none;
    }
    .xterm .composition-view {
        /* TODO: Composition position got messed up somewhere */
        background: #000;
        color: #FFF;
        display: none;
        position: absolute;
        white-space: nowrap;
        z-index: 1;
    }
    .xterm .composition-view.active {
        display: block;
    }
    .xterm .xterm-viewport {
        /* On OS X this is required in order for the scroll bar to appear fully opaque */
        background-color: #000;
        overflow-y: scroll;
        cursor: default;
        position: absolute;
        right: 0;
        left: 0;
        top: 0;
        bottom: 0;
    }
    .xterm .xterm-screen {
        position: relative;
    }
    .xterm .xterm-screen canvas {
        position: absolute;
        left: 0;
        top: 0;
    }
    .xterm .xterm-scroll-area {
        visibility: hidden;
    }
    .xterm-char-measure-element {
        display: inline-block;
        visibility: hidden;
        position: absolute;
        top: 0;
        left: -9999em;
        line-height: normal;
    }
    .xterm {
        cursor: text;
    }
    .xterm.enable-mouse-events {
        /* When mouse events are enabled (eg. tmux), revert to the standard pointer cursor */
        cursor: default;
    }
    .xterm.xterm-cursor-pointer {
        cursor: pointer;
    }
    .xterm.column-select.focus {
        /* Column selection mode */
        cursor: crosshair;
    }
    .xterm .xterm-accessibility,
    .xterm .xterm-message {
        position: absolute;
        left: 0;
        top: 0;
        bottom: 0;
        right: 0;
        z-index: 10;
        color: transparent;
    }
    .xterm .live-region {
        position: absolute;
        left: -9999px;
        width: 1px;
        height: 1px;
        overflow: hidden;
    }
    .xterm-dim {
        opacity: 0.5;
    }
    .xterm-underline {
        text-decoration: underline;
    }
</style>