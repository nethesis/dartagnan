var StorageService = {
  methods: {
    get(prop, parse=true) {
      return parse ? JSON.parse(localStorage.getItem(prop)) : localStorage.getItem(prop)
    },
    set(prop, object) {
      localStorage.setItem(prop, JSON.stringify(object))
    },
    delete(prop) {
      localStorage.removeItem(prop)
    },
    clear() {
      localStorage.clear()
    }
  }
};
export default StorageService;
