<template>
    <span>{{ formatted(value) }}</span>
</template>

<script>
    import moment from 'moment';
    export default {
        name: 'DateTime',
        props: {
            value: String
        },
        methods: {
            formatted(t) {
                return moment(t).format('YYYY-MM-DD HH:mm:ss');
            }
        }
    }
</script>

