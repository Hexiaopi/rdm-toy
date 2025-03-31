import { defineStore } from 'pinia';

export const useConnStore = defineStore('conn', {
    state: () => ({
        current: '',
        list: [],
    })
})