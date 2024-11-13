# Aplikasi Web Donasi

Ini adalah aplikasi web yang memungkinkan pengguna untuk mengirim donasi dan menampilkannya secara real-time menggunakan teknologi WebSocket. Proyek ini terdiri dari dua komponen utama:
1. **Donasi Pendonor (Halaman Donor)**: Tempat di mana pengguna dapat mengirimkan donasi bersama dengan pesan pribadi.
2. **Penerima Donasi (Halaman Penerima)**: Tempat di mana donasi yang diterima ditampilkan secara real-time.

## Fitur

- **Update Donasi Real-time**: Donasi diproses dan ditampilkan secara real-time menggunakan WebSocket.
- **Form Donasi Interaktif**: Pengguna dapat memasukkan jumlah donasi dan pesan opsional untuk penerima.
- **Auto-Scroll untuk Daftar Donasi**: Halaman penerima otomatis menggulir ke bawah untuk menampilkan donasi terbaru.
- **UI Responsif**: Desain yang sederhana, ramah untuk perangkat mobile, dan menarik secara visual.

## Struktur Proyek

Proyek ini terdiri dari dua file HTML utama:
- **Donasi Pendonor (index.html)**: Halaman tempat pengguna dapat mengirimkan donasi.
- **Penerima Donasi (receiver.html)**: Halaman yang menampilkan donasi yang masuk secara real-time.

### Donasi Pendonor (index.html)

Halaman ini berisi form input di mana pengguna dapat:
- Memasukkan jumlah donasi dalam IDR (Rupiah).
- Menambahkan pesan untuk menyertai donasi mereka (opsional).

Donasi akan dikirimkan ke server WebSocket dalam format JSON, dan setelah berhasil dikirim, pesan sukses akan ditampilkan kepada pengguna.
