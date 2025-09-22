# Cloud Native Go ☁️

This repo is where I’m documenting my journey of learning **Cloud Native development with Go**.  
I’m following along with the book *Cloud Native Go* by Matthew A. Titmus and experimenting with different patterns and examples.



### 🔍 What’s inside?

You’ll find small, runnable Go programs that explore concepts like:

- **Concurrency patterns** → debounce, retry, throttle, timeout  
- **Resilience patterns** → retries with backoff, cancellation with context  
- **Simple services** → an HTTP key-value store built with Gorilla Mux  

The idea is to keep each pattern/service as a small, focused example you can run and play around with.



### 🛠 Running the code

Clone the repo:

```bash
git clone https://github.com/Samarthasbhat/Cloud-Native-Go.git
cd Cloud-Native-Go

```
### Project Layout 📂

```
Cloud-Native-Go/
├── concurrency-patterns/
│   ├── debounce/
│   ├── retry/
│   ├── throttle/
│   └── timeout/
├── cloud-native-service/
│   └── key-value-store/
└── README.md
```

### 🚀 Contributing 
This is mostly a personal learning repo, but if you have improvements, new patterns, or better examples, feel free to:

- Fork the repo
- Add your changes
- Open a pull request

I’d love to learn from other approaches too.

### 📚 References 

- Cloud Native Go — Matthew A. Titmus
- Go docs → context, net/http
- Gorilla Mux
