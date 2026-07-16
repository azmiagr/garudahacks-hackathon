# PRODUCT REQUIREMENTS DOCUMENT

# PijarNusa
### Platform Ekosistem Logistik Kebencanaan Terpadu

*Peta Geospasial Real-Time • Order Logistik Otonom • Rantai Kustodi Terverifikasi*

**Versi 1.0 | 16 Juli 2026**
**Status:** Draft untuk Review — Disusun berdasarkan Konsep Solusi GarudaHack 7.0
**Disusun oleh:** Tim Produk & Rekayasa PijarNusa

---

## Kontrol Dokumen

| Versi | Tanggal | Penulis | Ringkasan Perubahan | Status |
|---|---|---|---|---|
| 0.1 | 10 Jul 2026 | Tim Hackathon PijarNusa | Draf awal: Konsep solusi, VPC, User Journey, arsitektur teknis dasar | Draft |
| 1.0 | 16 Jul 2026 | Tim Produk & Rekayasa | PRD lengkap: persona, functional & non-functional requirements, data model, API, aturan bisnis, edge case, RICE, roadmap | Review |

### Daftar Distribusi & Persetujuan

| Peran | Nama/Fungsi | Tanggung Jawab Review |
|---|---|---|
| Product Owner | Ketua Tim Produk | Kelayakan bisnis, prioritas fitur, ROI |
| Tech Lead | Ketua Tim Rekayasa | Kelayakan teknis, arsitektur, estimasi effort |
| Design Lead | UX/UI Researcher | Kelayakan alur pengguna & aksesibilitas |
| Trust & Safety | Compliance Officer | Kepatuhan UU PDP, regulasi penggalangan dana, anti-fraud |
| Mitra Eksternal | Payment Gateway (Xendit/Midtrans) | Kelayakan disbursement & KYC merchant |

---

## Daftar Isi

1. [Ringkasan Eksekutif](#1-ringkasan-eksekutif)
2. [Latar Belakang & Pernyataan Masalah](#2-latar-belakang--pernyataan-masalah)
3. [Tujuan Produk & Metrik Keberhasilan](#3-tujuan-produk--metrik-keberhasilan)
4. [Persona Pengguna](#4-persona-pengguna)
5. [Ruang Lingkup Produk](#5-ruang-lingkup-produk)
6. [Peta Perjalanan Pengguna (User Journey Map)](#6-peta-perjalanan-pengguna-user-journey-map)
7. [Kebutuhan Fungsional (Functional Requirements)](#7-kebutuhan-fungsional-functional-requirements)
8. [Aturan Bisnis & Penanganan Kasus Tepi (Edge Cases)](#8-aturan-bisnis--penanganan-kasus-tepi-edge-cases)
9. [Kebutuhan Non-Fungsional (Non-Functional Requirements)](#9-kebutuhan-non-fungsional-non-functional-requirements)
10. [Arsitektur Teknis & Pipeline Sistem](#10-arsitektur-teknis--pipeline-sistem)
11. [Rencana Analitik & Instrumentasi](#11-rencana-analitik--instrumentasi)
12. [Risiko & Mitigasi](#12-risiko--mitigasi)
13. [Prioritisasi Fitur (Kerangka RICE)](#13-prioritisasi-fitur-kerangka-rice)
14. [Roadmap & Fase Peluncuran](#14-roadmap--fase-peluncuran)
15. [Asumsi & Pertanyaan Terbuka](#15-asumsi--pertanyaan-terbuka)
16. [Glosarium](#16-glosarium)

---

## 1. Ringkasan Eksekutif

PijarNusa adalah platform ekosistem logistik kebencanaan terpadu yang menghubungkan lima aktor — Admin Posko, Donatur, Toko Mitra, Relawan Kurir, dan Penyintas — dalam satu siklus tertutup yang dapat diaudit sepenuhnya, mulai dari donasi masuk hingga barang diterima di lokasi bencana. Produk ini mengganti model donasi konvensional yang buram dengan tiga pilar inti: visualisasi geospasial real-time untuk memicu empati donatur, mekanisme order logistik otonom bergaya ride-hailing untuk mempercepat pemenuhan kebutuhan, dan rantai kustodi digital (chain of custody) berbasis QR dinamis serta bukti foto wajib untuk menjamin barang benar-benar sampai ke penyintas.

Dokumen ini menyempurnakan konsep solusi awal (hasil GarudaHack 7.0) menjadi PRD yang siap masuk fase pengembangan produksi. Selain merapikan lima modul inti yang sudah dirancang, dokumen ini menambahkan komponen yang sebelumnya belum tercakup namun krusial untuk peluncuran nyata: modul Trust & Safety / backoffice admin, alur dispute resolution dan refund, kebijakan dana mengendap (idle fund), fallback konektivitas rendah di zona bencana, model data lengkap, spesifikasi API, kerangka kepatuhan regulasi (UU PDP, izin penggalangan dana Kemensos, ketentuan Penyelenggara Jasa Pembayaran BI/OJK), serta rencana mitigasi risiko operasional dan reputasi.

**Masalah inti:** donatur tidak percaya donasinya digunakan tepat sasaran, posko bencana kesulitan mendapatkan logistik yang relevan dengan cepat, toko lokal enggan menyuplai barang dalam jumlah besar tanpa jaminan pembayaran, dan relawan kurir rentan dituduh kehilangan barang di jalan.

**Solusi:** satu platform tertutup yang mengunci dana donatur ke pesanan spesifik, mendistribusikan pesanan secara otomatis ke toko terdekat, mencatat setiap perpindahan barang dengan QR dinamis dan hash rantai (mock-ledger), serta memaksa bukti foto akhir sebelum dana dicairkan — menghasilkan transparansi end-to-end yang dapat diverifikasi publik.

**Target dampak (12 bulan pertama):** menjangkau 150 posko bencana aktif, memfasilitasi lebih dari Rp15 miliar dana tersalurkan, mempertahankan waktu rata-rata dari dana terkumpul hingga barang diterima di bawah 6 jam untuk area Jawa dan 24 jam untuk luar Jawa, serta mencapai skor kepercayaan donatur (net trust score) di atas 80%.

---

## 2. Latar Belakang & Pernyataan Masalah

### 2.1 Konteks Pasar

Indonesia adalah salah satu negara paling rawan bencana di dunia. Sepanjang tahun 2025, **Badan Nasional Penanggulangan Bencana (BNPB) mencatat 3.233 kejadian bencana alam**, didominasi banjir (1.652 kejadian, 51,1%), cuaca ekstrem (714), kebakaran hutan dan lahan (546), serta tanah longsor (233). Rangkaian bencana ini mengakibatkan sekitar 1.623 korban jiwa, 220 orang hilang, 5.713 luka-luka, dan lebih dari 10 juta warga terdampak atau mengungsi, dengan 220.560 bangunan rusak — 98,2% di antaranya rumah tinggal.

*Sumber: Rekapitulasi BNPB 2025, dipublikasikan melalui indonesiabaik.id dan Portal Satu Data Bencana Indonesia (data.bnpb.go.id).*

Skala kejadian sebesar ini menciptakan permintaan logistik darurat yang sangat sering dan tersebar geografis, namun penyaluran bantuan di lapangan masih sangat bergantung pada koordinasi manual lintas relawan, grup WhatsApp, dan rekening pribadi — pola yang secara struktural sulit diaudit dan lambat merespons kebutuhan riil di posko.

### 2.2 Masalah yang Divalidasi per Aktor

| Aktor | Pain Point Utama | Konsekuensi Bisnis |
|---|---|---|
| Donatur | Donasi terasa seperti "kotak hitam" tanpa laporan penggunaan dana yang jelas; tidak ada bukti barang benar-benar sampai | Donor fatigue, penurunan retensi donasi berulang, reputasi platform donasi menurun |
| Admin Posko | Bantuan datang lambat dan sering tidak sesuai kebutuhan aktual (mis. dikirim mi instan padahal butuh air bersih); pembukuan manual memakan waktu | Kebutuhan darurat penyintas tidak terpenuhi tepat waktu, admin posko rawan dituduh menyalahgunakan bantuan |
| Toko Mitra | Rawan diutangi relawan tanpa kepastian pembayaran saat kondisi krisis | Toko enggan menyuplai stok besar, pasokan logistik darurat tersendat |
| Relawan Kurir | Rawan dituduh kehilangan/menukar barang di jalan tanpa bukti serah terima yang sah | Menurunkan minat relawan berpartisipasi, memperlambat pengiriman |
| Penyintas | Tidak punya suara dalam proses; kerap menerima barang yang tidak relevan dengan kebutuhan riil | Distribusi bantuan tidak tepat sasaran, kesenjangan kebutuhan vs pasokan |

### 2.3 Kesenjangan pada Konsep Solusi Awal

Konsep solusi awal (dokumen GarudaHack 7.0) sudah kuat pada mekanisme inti — peta interaktif, bidding toko, dan rantai kustodi QR. Namun untuk siap produksi, PRD ini mengidentifikasi tujuh kesenjangan yang perlu ditutup:

1. Tidak ada peran moderasi/administrator platform (Trust & Safety) untuk menangani laporan kecurangan, sengketa, dan penonaktifan aktor bermasalah.
2. Tidak ada alur refund/dispute resolution ketika dana sudah terkumpul namun pesanan gagal, batal, atau terbukti fraud.
3. Tidak ada kebijakan eksplisit untuk dana lebih (overfunding) atau dana kurang (underfunding) pada satu posko.
4. Tidak ada mekanisme fallback saat konektivitas internet hilang total di lokasi bencana (bukan hanya kamera rusak).
5. Belum ada kerangka kepatuhan regulasi: UU Pelindungan Data Pribadi (UU PDP No. 27/2022), izin Pengumpulan Uang atau Barang (PUB) dari Kementerian Sosial, serta status penyelenggara pembayaran di bawah pengawasan Bank Indonesia/OJK.
6. Model data, kontrak API, dan definisi non-functional requirement (keamanan, performa, keandalan) belum dispesifikasikan.
7. Belum ada kerangka prioritisasi fitur (RICE) dan roadmap bertahap dari purwarupa hackathon menuju produk produksi.

Seluruh kesenjangan di atas dijawab secara rinci pada Bagian 7 hingga 14 dokumen ini.

---

## 3. Tujuan Produk & Metrik Keberhasilan

### 3.1 North Star Metric

**Rasio Penyelesaian Terverifikasi (Verified Fulfillment Rate)** — persentase permintaan logistik posko yang selesai penuh dengan bukti foto tervalidasi, hash rantai kustodi utuh, dan dana tersalurkan ke Toko Mitra, dihitung dari total permintaan yang didanai penuh. Metrik ini dipilih karena secara simultan merepresentasikan kepercayaan donatur, efisiensi operasional, dan integritas rantai kustodi.

### 3.2 Tujuan Bisnis & Target 12 Bulan

| Tujuan | Metrik | Baseline (Peluncuran) | Target 6 Bulan | Target 12 Bulan |
|---|---|---|---|---|
| Skala jangkauan | Jumlah posko aktif terlayani | 0 | 50 posko | 150 posko |
| Kepercayaan donatur | Net Trust Score (survei pasca-donasi) | N/A | 70% | 85% |
| Kecepatan respons | Waktu dana terkumpul → barang diterima (median) | N/A | 12 jam (Jawa) / 36 jam (luar Jawa) | 6 jam (Jawa) / 24 jam (luar Jawa) |
| Partisipasi Toko Mitra | Toko Mitra aktif terverifikasi | 0 | 300 toko | 1.200 toko |
| Retensi donatur | Donatur yang berdonasi ≥2x dalam 90 hari | N/A | 25% | 40% |
| Integritas rantai kustodi | Verified Fulfillment Rate (North Star) | N/A | 80% | 92% |
| Keberlanjutan finansial toko | Waktu pencairan dana ke Toko Mitra (median) | N/A | < 2 jam | < 30 menit |

### 3.3 Metrik Kesehatan Operasional (Guardrail Metrics)

- Tingkat sengketa (dispute rate) per pesanan selesai < 3%.
- Tingkat kegagalan pemindaian QR yang memerlukan fallback PIN < 8% dari total handshake.
- Tingkat penolakan foto verifikasi akhir (karena gelap/tidak relevan) < 5%.
- Uptime layanan inti (peta, donasi, order matching) ≥ 99,5% di luar masa insiden bencana skala nasional, ≥ 99,9% saat siaga darurat aktif.
- Tidak ada kebocoran data pribadi (PII) donatur, penyintas, atau mitra sepanjang tahun berjalan (zero-incident target).

Metrik guardrail dipantau berdampingan dengan metrik pertumbuhan agar tim tidak mengoptimalkan kecepatan/skala dengan mengorbankan integritas dan keamanan.

---

## 4. Persona Pengguna

Lima persona berikut disusun berdasarkan Value Proposition Canvas pada konsep awal, diperluas dengan demografi, konteks penggunaan perangkat, tingkat literasi digital, dan indikator keberhasilan spesifik untuk memandu keputusan desain dan prioritisasi fitur.

### Persona 1 — "Admin Posko Rian"
*Inisiator kebutuhan & penerima akhir logistik*

| Dimensi | Detail |
|---|---|
| Peran | Relawan garis depan / koordinator posko darurat (BPBD lokal, organisasi kemanusiaan, atau warga terdampak yang ditunjuk) |
| Usia & Konteks | 25–45 tahun, mengoperasikan aplikasi dari smartphone kelas menengah dengan sinyal 3G/4G tidak stabil di lokasi bencana |
| Jobs to be Done | Mengabarkan kondisi darurat secepat mungkin, mendapatkan logistik yang sesuai kebutuhan riil, mendistribusikan barang ke penyintas, mempertanggungjawabkan penggunaan bantuan |
| Pain Points | Bantuan tidak sesuai kebutuhan; pembukuan manual menyita waktu di tengah krisis; rawan dituduh menyalahgunakan donasi karena tidak ada bukti digital |
| Gains yang Dicari | Logistik tiba dalam hitungan jam; laporan otomatis yang membebaskannya dari tuduhan korupsi karena tidak pernah memegang uang tunai |
| Indikator Sukses | Waktu dari permintaan dibuat hingga logistik terdana penuh < 3 jam; 0 keluhan administratif terkait pertanggungjawaban dana |

### Persona 2 — "Donatur Maya"
*Individu penyumbang yang tergerak visual & butuh bukti dampak*

| Dimensi | Detail |
|---|---|
| Peran | Profesional urban usia 22–40 tahun yang aktif bermedia sosial, atau tim CSR korporasi yang butuh laporan audit |
| Konteks Penggunaan | Mengakses melalui smartphone di sela waktu luang; mengambil keputusan donasi dalam hitungan detik berdasarkan visual |
| Jobs to be Done | Menyalurkan dana secara efisien dan tepat sasaran, melacak penggunaan setiap rupiah, mendapatkan bukti nyata untuk kepuasan pribadi atau laporan audit korporasi |
| Pain Points | Donasi terasa seperti kotak hitam; bosan platform yang hanya meminta uang tanpa nilai balik; tidak yakin dana dibelikan barang yang relevan |
| Gains yang Dicari | Log transaksi terperinci, foto penyintas menerima barang yang ia danai, reward/insentif atas kebaikannya |
| Indikator Sukses | ≥70% donatur membuka dashboard log transaksinya dalam 7 hari; ≥40% donatur berdonasi ulang dalam 90 hari |

### Persona 3 — "Toko Mitra Pak Herman"
*Pemilik toko/grosir lokal, penyedia barang*

| Dimensi | Detail |
|---|---|
| Peran | Pemilik toko sembako, grosir, atau apotek lokal di sekitar radius bencana yang ingin memperoleh omzet tambahan |
| Konteks Penggunaan | Mengoperasikan aplikasi di sela melayani pembeli fisik; sensitif terhadap kepastian pembayaran |
| Jobs to be Done | Menjual stok barang dalam jumlah besar secara instan dengan risiko minimal |
| Pain Points | Rawan diutangi relawan tanpa kepastian bayar, terutama saat kondisi krisis membuat penagihan sulit dilakukan |
| Gains yang Dicari | Omzet besar dengan pembayaran dijamin sistem; validasi reputasi bisnis yang bisa dipakai untuk keperluan lain |
| Indikator Sukses | Waktu pencairan dana < 30 menit setelah verifikasi tiba; tingkat pembatalan pesanan oleh toko < 5% |

### Persona 4 — "Relawan Kurir Dewi"
*Individu/komunitas pengantar barang dari toko ke posko*

| Dimensi | Detail |
|---|---|
| Peran | Relawan independen atau anggota komunitas motor/mobil yang bersedia mengantar barang logistik |
| Konteks Penggunaan | Berkendara di area krisis dengan akses sinyal tidak menentu; membutuhkan proses handshake yang cepat dan tidak ribet |
| Jobs to be Done | Membantu mengantarkan barang logistik dengan aman dan mendapat pengakuan atas kontribusinya |
| Pain Points | Rawan dituduh menghilangkan atau menukar barang di jalan tanpa bukti serah terima yang sah |
| Gains yang Dicari | Validasi aktivitas sosial yang diakui platform, portofolio kontribusi yang bisa ditunjukkan |
| Indikator Sukses | Tingkat fallback PIN akibat gagal scan QR < 8%; 0 sengketa kehilangan barang yang tidak terselesaikan |

### Persona 5 — "Penyintas Bu Sari"
*Warga terdampak, penerima manfaat akhir*

| Dimensi | Detail |
|---|---|
| Peran | Warga terdampak bencana yang menjadi subjek bukti penerimaan akhir dan penerima manfaat logistik |
| Konteks Penggunaan | Umumnya tidak mengoperasikan aplikasi secara langsung; berinteraksi melalui Admin Posko sebagai perantara |
| Jobs to be Done | Mendapatkan kebutuhan dasar (air, makanan, popok, obat) secepat dan setepat mungkin |
| Pain Points | Sering menerima bantuan yang tidak relevan dengan kebutuhan riil karena kurangnya jalur pelaporan kebutuhan yang akurat |
| Gains yang Dicari | Kebutuhan darurat terpenuhi tanpa harus menunggu birokrasi panjang atau merasa menjadi objek eksploitasi foto |
| Indikator Sukses | Tingkat kesesuaian barang diterima vs kebutuhan yang dilaporkan Admin Posko ≥ 90% (diukur via survei sampling pasca-distribusi) |

### Persona Tambahan — "Admin Trust & Safety Nadia" (peran baru, lihat Bagian 7.10)

Staf internal PijarNusa yang memoderasi laporan kecurangan, memverifikasi KYC Toko Mitra dan Kurir, menangani sengketa, serta memiliki kewenangan membekukan dana atau menonaktifkan akun. Peran ini tidak ada pada konsep awal namun wajib ada sebelum peluncuran publik karena platform menangani dana pihak ketiga secara langsung.

---

## 5. Ruang Lingkup Produk

### 5.1 Dalam Lingkup — Fase MVP (Rilis Publik Pertama)

- Peta Bencana Interaktif dengan gradasi warna urgensi pendanaan per posko.
- Pembuatan & manajemen posko oleh Admin Posko (foto, daftar kebutuhan, radius geofencing).
- Deposit saldo donatur via QRIS/Virtual Account dan alokasi dana proporsional otomatis.
- Broadcast order ke Toko Mitra dengan mekanisme rebutan (first-accept-wins).
- Rantai kustodi dua tahap (Toko→Kurir, Kurir→Posko) dengan QR dinamis dan fallback PIN 6 digit.
- Verifikasi akhir wajib kamera (forced camera) dengan validasi geofencing GPS.
- Beranda Transparansi Publik & dashboard log transaksi per donatur.
- Disbursement otomatis ke rekening Toko Mitra via payment gateway.
- Program loyalitas dasar (poin yang dapat ditukar pulsa/voucer).
- Backoffice Trust & Safety untuk moderasi laporan, KYC, dan pembekuan dana (Bagian 7.10).
- Alur dispute resolution & refund dasar (Bagian 7.12).
- Notifikasi push dengan fallback SMS untuk area minim sinyal data.

### 5.2 Dalam Lingkup — Fase 2 (3–6 Bulan Pasca-Peluncuran)

- Mode offline penuh dengan sinkronisasi tertunda (offline-first PWA) untuk zona blackout sinyal total.
- Multi-item marketplace: donatur dapat memilih kebutuhan spesifik dalam satu posko (bukan hanya dana gabungan).
- Program reputasi lanjutan untuk Toko Mitra & Kurir (badge, prioritas order).
- Integrasi dengan BNPB/BPBD untuk validasi status posko resmi.
- Dukungan multi-bahasa daerah untuk Admin Posko di luar Jawa.

### 5.3 Dalam Lingkup — Fase 3 (6–12 Bulan)

- Prediksi kebutuhan logistik berbasis riwayat bencana serupa (model prediktif sederhana).
- Marketplace korporasi untuk donasi CSR terjadwal dan laporan audit otomatis.
- Ekspansi ke bantuan non-logistik: tenaga medis, tempat tinggal sementara (shelter matching).

### 5.4 Di Luar Lingkup (Eksplisit Tidak Dikerjakan)

- PijarNusa tidak menjadi penyedia layanan pengiriman barang komersial (bukan kompetitor logistik reguler); kurir murni berbasis kerelawanan bencana.
- Tidak menangani distribusi bantuan tunai langsung ke penyintas (cash transfer) pada fase MVP — hanya barang.
- Tidak menyediakan asuransi jiwa/kecelakaan untuk relawan kurir pada MVP (dicatat sebagai risiko pada Bagian 12).
- Tidak melakukan verifikasi lapangan fisik independen atas kondisi bencana — bergantung pada laporan Admin Posko terverifikasi dan pelaporan komunitas (flagging).
- Tidak mendukung donasi barang fisik langsung dari donatur (hanya dana yang dikonversi sistem menjadi pesanan barang ke Toko Mitra).

---

## 6. Peta Perjalanan Pengguna (User Journey Map)

Perjalanan berikut merangkum siklus tertutup lintas lima aktor, disempurnakan dengan detail respons sistem dan penanganan kegagalan pada setiap fase agar tim desain dan rekayasa memiliki acuan perilaku yang presisi, bukan hanya alur ideal (happy path).

### Fase 1: Discovery & Pendanaan Otonom

**Aksi Pengguna:**
- Admin Posko membuka aplikasi, memotret kondisi bencana, menulis daftar kebutuhan darurat dengan kuantitas spesifik, dan mengunggahnya beserta titik GPS posko.
- Sistem melakukan pengecekan duplikasi posko dalam radius 500m untuk mencegah entri ganda pada bencana yang sama.
- Donatur membuka Peta Bencana Interaktif, melihat titik merah berkedip bergradasi urgensi, mengklik titik untuk melihat rincian, lalu melakukan top-up saldo via QRIS/VA.
- Donatur dapat mengatur filter preferensi (jenis bencana, wilayah) agar dana otomatis dialokasikan ke posko yang cocok saat top-up berikutnya (mode "donasi otonom").

**Touchpoint:** Form Pembuatan Posko, Modul Kamera, Peta Interaktif, Gateway Pembayaran QRIS/VA, Modul Filter Preferensi Donatur
**Emosi:** Admin Posko: kritis dan mendesak. Donatur: tersentuh secara visual, terdorong bertindak cepat.
**Pain Point:** Donatur bingung memilih posko yang paling membutuhkan; Admin Posko cemas menunggu respons.
**Antisipasi/Respons Sistem:** Peta menggunakan gradasi warna bahaya (merah tua = belum terdanai sama sekali) sehingga atensi otomatis tertuju ke titik paling darurat. Begitu dana mencukupi, status posko berubah instan menjadi "Terdanai Penuh", mengubah kecemasan Admin Posko menjadi kelegaan langsung.

### Fase 2: Order Matching & Persiapan Toko

**Aksi Pengguna:**
- Sistem memancarkan (broadcast) pesanan ke seluruh Toko Mitra terverifikasi dalam radius konfigurasi (default 5 km, membesar bertahap jika tidak ada respons).
- Toko Mitra menekan "Setujui & Siapkan"; toko lain otomatis kehilangan akses ke pesanan tersebut (locking pesimistik di level database).
- Layar toko menampilkan spanduk konfirmasi dana terkunci beserta rincian nominal.

**Touchpoint:** Dashboard Toko Mitra, Modul Order Bidding, Notifikasi Instan/SMS
**Emosi:** Toko antusias mendapat omzet di masa krisis, namun tetap waspada terhadap kepastian bayar.
**Pain Point:** Toko ragu menyiapkan barang banyak bila tidak ada jaminan tertulis; risiko tidak ada toko yang merespons dalam radius awal.
**Antisipasi/Respons Sistem:** Dana ditandai "locked" secara atomik saat broadcast dikirim (bukan saat toko menyetujui) sehingga tidak ada race condition dana terpakai dua kali. Jika tidak ada toko merespons dalam 10 menit, radius broadcast otomatis melebar 2x dan notifikasi eskalasi dikirim ke Admin Trust & Safety (lihat Bagian 8.2).

### Fase 3: Rantai Kustodi Berjenjang

**Aksi Pengguna:**
- Handshake 1: Kurir tiba di toko, memindai QR dinamis milik toko (berganti tiap 30 detik); kustodi barang berpindah resmi ke Kurir.
- Handshake 2: Kurir tiba di posko, menampilkan QR dinamis miliknya; Admin Posko memindai; kustodi berpindah resmi ke Posko.
- Bila kamera atau layar rusak, aktor menggunakan fallback PIN 6 digit yang juga berganti tiap 30 detik, mengikuti pola token OTP.

**Touchpoint:** Modul QR Generator Dinamis, Pemindai QR In-App, Modul Fallback PIN
**Emosi:** Rasa tanggung jawab tinggi atas barang yang dibawa; kurir merasa tegang membawa barang berharga di area krisis.
**Pain Point:** Kegagalan pemindaian akibat kamera rusak, layar retak, atau koneksi terputus saat validasi ke backend.
**Antisipasi/Respons Sistem:** Setiap pemindaian yang tervalidasi langsung memindahkan tanggung jawab digital secara resmi, membebaskan kurir dari tuduhan penggelapan. QR dan PIN dibuat dengan validitas offline-tolerant: token dibangkitkan di backend namun di-cache di klien selama 90 detik agar tetap dapat divalidasi walau ada keterlambatan sinkronisasi jaringan (lihat Bagian 8.7).

### Fase 4: Verifikasi Akhir, Transparansi, & Retensi Loyalitas

**Aksi Pengguna:**
- Admin Posko menekan "Selesaikan Pesanan"; aplikasi memaksa kamera terbuka (galeri terkunci) dan mengekstraksi koordinat GPS serta timestamp UTC.
- Backend memvalidasi geofencing (jarak GPS akhir vs GPS posko awal < 500m); bila valid, ledger dikunci dan disbursement dipicu.
- Sistem mempublikasikan foto ke Beranda Transparansi Publik dengan overlay nominal, jumlah donatur, dan daftar barang.
- Donatur, Toko, dan Kurir menerima notifikasi dan dapat mengklaim poin loyalitas di dashboard masing-masing.

**Touchpoint:** Modul Forced Camera, Beranda Publik, Dashboard Transaksi Donatur, Dashboard Analitik Toko & Kurir, Modul Redeem Reward
**Emosi:** Kepuasan batin, kebanggaan publik, dan keinginan donatur menyumbang lagi untuk mengejar level reward.
**Pain Point:** Admin Posko malas atau tergesa mengambil foto akhir yang tidak relevan/gelap; validasi geofencing gagal karena GPS lemah di dalam bangunan.
**Antisipasi/Respons Sistem:** Foto yang terindikasi gelap/buram otomatis ditolak sistem (deteksi kecerahan minimum) dan diminta pengulangan; setelah 2x penolakan otomatis, kasus dieskalasi ke Admin Trust & Safety untuk verifikasi manual. Laporan publik (flagging) atas foto yang tidak relevan menurunkan skor reputasi posko. Kegagalan geofencing akibat GPS lemah memicu verifikasi sekunder berbasis triangulasi menara seluler sebelum eskalasi manual — lihat Bagian 8.4.

---

## 7. Kebutuhan Fungsional (Functional Requirements)

Setiap modul di bawah ini mengikuti format Epic → User Story → Acceptance Criteria agar dapat langsung diterjemahkan menjadi backlog pengembangan. Modul 7.1, 7.10, 7.11, dan 7.12 merupakan penambahan baru terhadap konsep awal untuk menutup kesenjangan kesiapan produksi.

### 7.1 Modul Autentikasi, Onboarding & KYC

Menangani pendaftaran dan verifikasi identitas seluruh aktor. Toko Mitra dan Relawan Kurir wajib melalui proses Know-Your-Customer (KYC) karena keduanya menerima/menyalurkan dana atau barang bernilai ekonomi, sementara Donatur dan Admin Posko menggunakan verifikasi ringan (lightweight) agar tidak menjadi friksi saat urgensi tinggi.

**User Stories:**
- Sebagai Toko Mitra, saya perlu mendaftarkan NIB/NPWP toko, foto KTP pemilik, dan nomor rekening bank agar dapat menerima disbursement.
- Sebagai Relawan Kurir, saya perlu memverifikasi NIK dan nomor HP aktif agar identitas saya tercatat dalam rantai kustodi.
- Sebagai Donatur, saya cukup mendaftar dengan nomor HP/email dan verifikasi OTP agar proses donasi tidak berbelit di momen genting.
- Sebagai Admin Posko, saya perlu diverifikasi manual oleh Trust & Safety (atau melalui rekomendasi BPBD/organisasi mitra terdaftar) sebelum dapat membuat posko pertama, untuk mencegah posko palsu.

**Kriteria Penerimaan (Acceptance Criteria):**
1. Registrasi Toko Mitra tidak dapat menerima order pertama sebelum status KYC = "Terverifikasi" oleh sistem otomatis (validasi NIB via API OSS) atau manual oleh Trust & Safety maksimal 1x24 jam.
2. Registrasi Kurir baru dibatasi menerima maksimal 1 order aktif dalam 48 jam pertama (masa percobaan) untuk mitigasi risiko sebelum reputasi terbentuk.
3. Admin Posko pertama kali wajib mengunggah surat tugas/afiliasi organisasi ATAU mendapat vouching dari Admin Posko lain yang sudah terverifikasi ≥3 posko sukses.
4. Seluruh data KYC dienkripsi at-rest dan hanya dapat diakses oleh peran Trust & Safety dengan audit log setiap kali dibuka (lihat Bagian 9.1).

### 7.2 Modul Peta Bencana Interaktif & Manajemen Posko

Beranda utama aplikasi berbasis peta yang menampilkan seluruh posko aktif dengan visual gradasi urgensi pendanaan, menjadi titik pemicu emosional bagi donatur sekaligus kanal input kebutuhan bagi Admin Posko.

**User Stories:**
- Sebagai Donatur, saya ingin melihat titik-titik posko pada peta dengan warna yang menunjukkan seberapa mendesak kebutuhan pendanaannya.
- Sebagai Admin Posko, saya ingin membuat entri posko baru dengan foto, daftar kebutuhan terstruktur (nama barang, satuan, kuantitas), dan radius geofencing dalam waktu kurang dari 3 menit.
- Sebagai Donatur, saya ingin memfilter peta berdasarkan jenis bencana dan wilayah agar dapat fokus pada isu yang saya pedulikan.

**Kriteria Penerimaan (Acceptance Criteria):**
1. Warna titik posko dihitung dari rasio dana_terkumpul/dana_dibutuhkan: merah tua (0–20%), merah (21–50%), oranye (51–80%), kuning (81–99%), hijau (100%, posko tidak lagi menerima donasi baru untuk kebutuhan yang sama).
2. Sistem menolak pembuatan posko baru dalam radius 500 meter dari posko aktif lain dengan kategori bencana yang sama dalam 72 jam terakhir; menampilkan opsi "gabung ke posko existing" sebagai gantinya.
3. Peta menggunakan clustering marker saat >50 posko berada dalam satu viewport untuk menjaga performa render di perangkat kelas menengah.
4. Setiap daftar kebutuhan wajib memiliki minimal 1 item dengan kuantitas dan estimasi harga satuan (diambil dari katalog referensi harga pasar) agar sistem dapat menghitung target dana otomatis.

### 7.3 Modul Donasi & Wallet

Mengelola alur dana masuk dari donatur, saldo mengendap (idle balance), dan alokasi proporsional dana ke pesanan logistik spesifik.

**User Stories:**
- Sebagai Donatur, saya ingin top-up saldo via QRIS atau Virtual Account dan langsung melihatnya teralokasi ke posko pilihan saya.
- Sebagai Donatur, saya ingin mengaktifkan mode "donasi otonom" agar saldo saya otomatis dialokasikan ke posko yang cocok dengan preferensi tanpa perlu memilih manual setiap saat.
- Sebagai Product Owner, saya perlu kebijakan yang jelas untuk dana lebih (overfunding) dan dana kurang (underfunding) pada satu posko agar tidak ada dana donatur yang "menggantung" tanpa kejelasan.

**Kriteria Penerimaan (Acceptance Criteria):**
1. Alokasi dana ke posko bersifat proporsional real-time: begitu donatur menekan "Donasi", saldo langsung didebit dan dicatat sebagai kontribusi terkunci (locked) ke posko tersebut — tidak ada status "pending" yang menggantung lebih dari beberapa detik.
2. Overfunding (dana terkumpul melebihi kebutuhan): kelebihan dana otomatis dialihkan ke "Dana Cadangan Posko" yang dapat digunakan Admin Posko untuk kebutuhan susulan pada posko yang sama, dengan notifikasi transparan ke seluruh donatur yang berkontribusi.
3. Underfunding (posko tidak mencapai target dana dalam 7 hari, umumnya karena kebutuhan sudah tidak relevan): dana proporsional dikembalikan ke saldo masing-masing donatur secara otomatis, dengan opsi "alihkan ke posko lain" satu klik.
4. Seluruh transaksi wallet (top-up, alokasi, refund, penarikan poin) tercatat sebagai entri di Merklized Mock-Ledger (lihat Bagian 10.3) sehingga riwayat tidak dapat diubah retroaktif.

### 7.4 Modul Order Matching & Bidding Toko

Mendistribusikan pesanan logistik yang sudah terdanai penuh ke Toko Mitra terdekat menggunakan mekanisme rebutan mirip aplikasi ride-hailing.

**User Stories:**
- Sebagai Toko Mitra, saya ingin menerima notifikasi instan saat ada pesanan yang cocok dengan kategori toko saya dan lokasi saya.
- Sebagai Toko Mitra, saya ingin memiliki jendela waktu yang adil untuk menekan "Setujui" tanpa dikalahkan oleh latensi jaringan yang tidak saya kendalikan.
- Sebagai Admin Posko, saya ingin tahu jika tidak ada toko yang merespons pesanan saya agar dapat mengambil tindakan lanjutan.

**Kriteria Penerimaan (Acceptance Criteria):**
1. Broadcast awal dikirim ke seluruh Toko Mitra terverifikasi dalam radius 5 km dengan kategori barang yang cocok; jika tidak ada respons dalam 10 menit, radius otomatis membesar ke 10 km, lalu 20 km pada menit ke-20.
2. Locking pesanan menggunakan mekanisme atomic compare-and-swap di backend (bukan first-request-wins di sisi klien) untuk menghindari dua toko sama-sama mendapat konfirmasi karena race condition jaringan.
3. Toko yang membatalkan pesanan setelah menekan "Setujui" (sebelum handshake pertama) dikenai penurunan skor reputasi dan pesanan otomatis di-broadcast ulang ke toko berikutnya tanpa perlu campur tangan Admin Posko.
4. Jika tidak ada toko merespons setelah radius maksimum (20 km) dan 30 menit berlalu, sistem membuat tiket eskalasi otomatis ke Admin Trust & Safety untuk pencarian toko manual atau redistribusi dana.

### 7.5 Modul Rantai Kustodi (Chain of Custody)

Mengawal perpindahan fisik barang dari Toko ke Kurir dan dari Kurir ke Posko menggunakan QR Code dinamis dengan fallback PIN, memastikan setiap perpindahan tercatat dan tidak dapat disangkal (non-repudiation).

**User Stories:**
- Sebagai Relawan Kurir, saya ingin memindai QR sederhana untuk mengonfirmasi saya benar-benar mengambil barang dari toko yang tepat.
- Sebagai Admin Posko, saya ingin memindai QR milik kurir untuk mengonfirmasi barang yang tiba adalah barang yang dikirim, bukan tertukar.
- Sebagai Relawan Kurir, saya perlu cara alternatif mengonfirmasi serah terima jika kamera ponsel saya atau ponsel lawan transaksi rusak.

**Kriteria Penerimaan (Acceptance Criteria):**
1. QR Code dan PIN 6-digit dibangkitkan ulang setiap 30 detik dan ditandatangani (signed) oleh backend menggunakan HMAC dengan kunci sesi unik per transaksi agar tidak dapat direplikasi/diprediksi.
2. Validasi handshake wajib menyertakan koordinat GPS pemindai; jika jarak pemindai dan lokasi tujuan (toko/posko) > 300 meter, sistem menampilkan peringatan namun tetap mengizinkan (untuk toleransi akurasi GPS dalam bangunan), dicatat sebagai flag anomali tingkat rendah untuk audit.
3. Fallback PIN memiliki tingkat keamanan setara QR: PIN yang sama tidak dapat digunakan dua kali dan kedaluwarsa dalam 30 detik, sama seperti mekanisme token OTP.
4. Setiap perpindahan kustodi berhasil memicu penulisan entri baru ke ledger dengan hash SHA-256 yang menggabungkan hash transaksi sebelumnya (rantai Merkle sederhana), sehingga histori tidak dapat diubah tanpa merusak seluruh rantai berikutnya.

### 7.6 Modul Verifikasi Akhir & Forced Camera

Mengunci siklus pesanan dengan bukti foto wajib di lokasi, memicu pencairan dana, dan mempublikasikan hasil ke Beranda Transparansi Publik.

**User Stories:**
- Sebagai Admin Posko, saya ingin proses verifikasi akhir sesederhana mungkin (satu tombol, satu foto) mengingat saya sedang sibuk menangani kondisi darurat.
- Sebagai Donatur, saya ingin yakin foto yang dipublikasikan benar-benar diambil di lokasi dan waktu yang sesuai, bukan foto lama atau dari galeri.

**Kriteria Penerimaan (Acceptance Criteria):**
1. Tombol "Selesaikan Pesanan" mengunci akses galeri perangkat sepenuhnya dan hanya membuka antarmuka kamera native/WebRTC — tidak ada jalur unggah file dari penyimpanan lokal.
2. Sistem mengekstraksi metadata EXIF dan HTML5 Geolocation API secara paksa saat rana kamera ditekan; foto tanpa metadata GPS valid otomatis ditolak dengan pesan error yang jelas.
3. Validasi geofencing membandingkan koordinat foto akhir dengan koordinat posko awal (toleransi radius 500 meter, dapat disesuaikan untuk posko dengan area luas seperti lapangan pengungsian).
4. Foto dengan skor kecerahan (brightness) di bawah ambang batas otomatis ditolak dan meminta pengulangan; setelah 2 kali penolakan berturut-turut, kasus dieskalasi otomatis untuk verifikasi manual Trust & Safety tanpa memblokir pencairan dana lebih dari 2 jam.
5. Begitu foto tervalidasi dan geofencing sesuai, ledger dikunci dengan status "Completed" secara atomik bersamaan dengan pemicuan webhook disbursement — tidak ada jeda antara status selesai dan pemicuan pembayaran.

### 7.7 Modul Transparansi Publik & Pelaporan Komunitas

Beranda publik yang menayangkan bukti penyelesaian setiap pesanan beserta mekanisme pelaporan (flagging) oleh masyarakat umum untuk menjaga integritas konten.

**User Stories:**
- Sebagai pengunjung publik (tanpa akun), saya ingin melihat bukti bahwa donasi di platform ini benar-benar sampai ke penyintas.
- Sebagai pengguna terdaftar, saya ingin dapat menandai (flag) posko yang fotonya terlihat tidak relevan, gelap, atau mencurigakan.

**Kriteria Penerimaan (Acceptance Criteria):**
1. Beranda publik dapat diakses tanpa login dan menampilkan foto, nominal, jumlah donatur, dan daftar barang untuk setiap pesanan berstatus "Completed".
2. Setiap entri publik menampilkan tautan "Lacak Rantai Kustodi" yang menunjukkan histori hash tanpa membuka data pribadi (nama disamarkan menjadi inisial, sesuai Bagian 9.1).
3. Fitur flag memerlukan alasan terstruktur (foto tidak jelas/barang tidak sesuai/dugaan fraud/lainnya); 3 flag independen dalam 24 jam otomatis menurunkan skor reputasi posko dan membekukan pencairan dana pesanan terkait hingga ditinjau Trust & Safety.

### 7.8 Modul Loyalitas & Reward

Program insentif gamifikasi untuk mendorong partisipasi berulang dari Donatur, Toko Mitra, dan Relawan Kurir.

**User Stories:**
- Sebagai Donatur, saya ingin mengumpulkan poin dari setiap donasi dan menukarnya dengan pulsa atau voucer e-commerce.
- Sebagai Toko Mitra dan Relawan Kurir, saya ingin memiliki dashboard "Jejak Kebaikan" yang merekap total kontribusi saya sebagai portofolio reputasi.

**Kriteria Penerimaan (Acceptance Criteria):**
1. Poin loyalitas didistribusikan otomatis dan atomik bersamaan dengan penguncian ledger "Completed" — tidak ada proses klaim manual terpisah yang dapat gagal.
2. Poin memiliki masa berlaku 12 bulan sejak diterbitkan dan ditampilkan dengan jelas tanggal kedaluwarsanya di dashboard.
3. Dashboard "Jejak Kebaikan" Toko Mitra dan Kurir dapat diekspor sebagai sertifikat PDF berisi ringkasan kontribusi terverifikasi (nominal, jumlah pesanan, periode) untuk keperluan reputasi bisnis atau administrasi kerelawanan.

### 7.9 Modul Disbursement & Pembayaran

Mengeksekusi pencairan dana fiat ke rekening Toko Mitra secara otomatis melalui payment gateway begitu verifikasi akhir selesai.

**User Stories:**
- Sebagai Toko Mitra, saya ingin dana langsung cair ke rekening bank saya tanpa proses manual tambahan setelah barang terverifikasi tiba.
- Sebagai Trust & Safety, saya perlu dapat menahan (hold) disbursement tertentu jika ada indikasi kecurangan sebelum dana benar-benar keluar dari platform.

**Kriteria Penerimaan (Acceptance Criteria):**
1. Webhook disbursement dipicu maksimal 5 detik setelah ledger status berubah menjadi "Completed"; kegagalan panggilan API payment gateway di-retry dengan exponential backoff hingga 5 kali sebelum dieskalasi ke antrian manual.
2. Setiap disbursement tercatat dengan idempotency key unik per pesanan untuk mencegah pencairan ganda meski terjadi retry.
3. Trust & Safety memiliki kewenangan mengaktifkan "pre-disbursement hold" maksimal 24 jam pada pesanan yang menerima ≥3 flag komunitas atau terdeteksi anomali geofencing, sebelum dana benar-benar ditransfer keluar platform.

### 7.10 Modul Backoffice Trust & Safety (Baru)

Dashboard administratif internal — tidak tercakup pada konsep awal — untuk staf PijarNusa memoderasi platform, memverifikasi KYC, menangani sengketa, dan menjaga integritas dana pihak ketiga. Modul ini wajib ada sebelum peluncuran publik karena platform memegang kepercayaan finansial banyak pihak.

**User Stories:**
- Sebagai Admin Trust & Safety, saya ingin antrean terpusat berisi seluruh laporan flag, eskalasi order-matching gagal, dan kasus verifikasi foto ditolak, terurut berdasarkan urgensi.
- Sebagai Admin Trust & Safety, saya ingin dapat membekukan akun (Toko/Kurir/Admin Posko) dan menahan dana terkait dalam satu tindakan saat kecurangan terbukti.
- Sebagai Admin Trust & Safety, saya ingin melihat riwayat lengkap rantai kustodi dan ledger hash untuk satu pesanan guna investigasi sengketa.

**Kriteria Penerimaan (Acceptance Criteria):**
1. Seluruh tindakan Trust & Safety (bekukan akun, tahan dana, setujui refund) memerlukan autentikasi dua faktor dan tercatat pada audit log yang tidak dapat dihapus.
2. Dashboard menampilkan SLA per kasus: eskalasi order-matching harus direspons < 30 menit, laporan flag prioritas tinggi (dugaan fraud) < 2 jam, verifikasi foto ditolak < 2 jam.
3. Peran dibagi minimal dua tingkat: "Reviewer" (dapat meninjau dan merekomendasikan) dan "Approver" (dapat mengeksekusi pembekuan dana/refund) — prinsip pemisahan tugas (segregation of duties) untuk mencegah penyalahgunaan wewenang internal.

### 7.11 Modul Notifikasi & Ketahanan Konektivitas Rendah (Baru)

Sistem notifikasi multi-kanal dengan fallback untuk area bencana yang minim sinyal data namun masih memiliki jaringan seluler dasar (2G/SMS).

**User Stories:**
- Sebagai Toko Mitra, saya ingin tetap menerima notifikasi pesanan baru meski koneksi data saya lemah.
- Sebagai Relawan Kurir, saya ingin QR/PIN saya tetap dapat divalidasi walau koneksi sempat terputus saat proses handshake.

**Kriteria Penerimaan (Acceptance Criteria):**
1. Notifikasi kritis (broadcast order baru, permintaan handshake) dikirim melalui push notification utama dan SMS fallback otomatis jika push tidak ter-deliver dalam 60 detik.
2. Aplikasi PWA melakukan caching lokal untuk QR/PIN aktif dan antrean aksi (scan, foto) yang disinkronkan otomatis ke backend begitu koneksi pulih, dengan indikator status "tersimpan lokal — menunggu sinkronisasi" yang jelas bagi pengguna.
3. Seluruh aksi yang disinkronkan ulang menyertakan timestamp asli saat aksi dilakukan (bukan timestamp saat sinkronisasi) untuk menjaga akurasi urutan rantai kustodi.

### 7.12 Modul Resolusi Sengketa & Refund (Baru)

Alur formal untuk menangani kegagalan pesanan, dugaan fraud, atau ketidaksesuaian barang — komponen yang wajib ada begitu platform menangani dana pihak ketiga dalam skala nyata.

**User Stories:**
- Sebagai Donatur, saya ingin mengajukan keberatan jika saya menemukan bukti kuat bahwa dana saya disalahgunakan.
- Sebagai Toko Mitra, saya ingin mengajukan sengketa jika pembayaran tertahan padahal saya sudah menyerahkan barang sesuai prosedur.
- Sebagai Relawan Kurir, saya ingin ada jalur pembelaan formal jika dituduh menghilangkan barang, dengan rantai kustodi sebagai bukti objektif.

**Kriteria Penerimaan (Acceptance Criteria):**
1. Setiap pesanan yang dieskalasi (via flag komunitas, pre-disbursement hold, atau pengajuan sengketa manual) mendapat status "Under Review" yang terlihat oleh seluruh aktor terkait, mencegah tindakan sepihak lebih lanjut selama investigasi.
2. Keputusan Trust & Safety atas sengketa wajib disertai rujukan eksplisit ke entri ledger/hash rantai kustodi yang relevan sebagai dasar keputusan (bukan opini subjektif semata).
3. Skenario refund parsial didukung: misalnya jika barang terbukti hilang di tengah rute (kustodi terakhir tercatat di Kurir, tidak pernah sampai ke Posko), dana yang belum dicairkan ke Toko dikembalikan proporsional ke donatur, sementara Toko tetap dibayar penuh karena kustodinya sudah berpindah sah ke Kurir — kerugian ditanggung melalui Dana Cadangan Platform (lihat Bagian 12, risiko operasional).
4. SLA penyelesaian sengketa: maksimal 3x24 jam untuk kasus standar, 1x24 jam untuk kasus yang melibatkan indikasi fraud aktif.

---

## 8. Aturan Bisnis & Penanganan Kasus Tepi (Edge Cases)

Bagian ini mendefinisikan secara eksplisit bagaimana sistem berperilaku di luar jalur ideal (happy path), karena justru pada kondisi darurat bencana-lah kasus tepi paling sering terjadi dan paling berisiko terhadap kepercayaan pengguna.

### 8.1 Dana Berlebih pada Satu Posko (Overfunding)
- Kelebihan dana di atas kebutuhan tercatat dialihkan otomatis ke "Dana Cadangan Posko", dapat digunakan Admin Posko untuk kebutuhan susulan pada bencana yang sama tanpa perlu donasi baru.
- Seluruh donatur yang berkontribusi menerima notifikasi transparan berisi rincian penggunaan dana cadangan tersebut begitu digunakan.
- Jika dalam 14 hari dana cadangan tidak digunakan, sistem menawarkan donatur opsi menarik kembali kontribusinya atau mengalihkannya ke posko aktif lain.

### 8.2 Tidak Ada Toko Merespons Order
- Radius broadcast melebar bertahap: 5 km (0 menit) → 10 km (10 menit) → 20 km (20 menit).
- Setelah 30 menit tanpa respons pada radius maksimum, tiket eskalasi otomatis dibuat ke Trust & Safety untuk pencarian manual atau kontak toko cadangan (mitra korporasi besar dengan cakupan nasional).
- Admin Posko menerima notifikasi status setiap tahap eskalasi agar tidak merasa permintaannya diabaikan.

### 8.3 Toko Membatalkan Setelah Menyetujui
- Pembatalan sebelum handshake pertama: skor reputasi toko turun, pesanan otomatis di-broadcast ulang ke toko berikutnya dalam radius yang sama.
- Pembatalan berulang (≥3 kali dalam 30 hari) memicu peninjauan status KYC toko oleh Trust & Safety dan dapat berujung suspensi sementara.

### 8.4 Kegagalan Validasi GPS/Geofencing
- Toleransi radius default 500 meter dari titik posko; dapat disesuaikan lebih luas untuk lokasi seperti lapangan pengungsian atau kompleks besar.
- Kegagalan geofencing pertama memicu verifikasi sekunder berbasis triangulasi menara seluler sebagai pembanding sebelum otomatis menolak.
- Kegagalan berulang (2x) mengunci pencairan dana sementara dan mengeskalasi ke Trust & Safety untuk verifikasi manual berbasis foto dan konteks lain (bukan GPS semata), dengan SLA maksimal 2 jam agar tidak menghambat operasional posko.

### 8.5 Barang Hilang di Tengah Rute Kustodi
- Rantai kustodi hash menunjukkan titik terakhir barang tercatat sah — ini menjadi bukti objektif utama investigasi, bukan testimoni sepihak.
- Jika kustodi terakhir tercatat di Kurir namun tidak pernah sampai ke Posko dalam 2x estimasi waktu tempuh normal, sistem otomatis membuka kasus "Barang Belum Tiba" dan menahan status penyelesaian.
- Toko tetap dibayar penuh (kustodinya sudah sah berpindah), sementara kerugian ditanggung Dana Cadangan Platform, dan kurir terkait menjalani peninjauan reputasi (lihat Bagian 12 untuk kebijakan pendanaan risiko ini).

### 8.6 Foto Verifikasi Akhir Ditolak Berulang
- Penolakan otomatis pertama (kecerahan/relevansi rendah): Admin Posko diminta mengulang foto dengan panduan visual (mis. "pastikan pencahayaan cukup, sertakan barang dan penyintas dalam satu bingkai").
- Penolakan kedua: kasus otomatis dieskalasi ke Trust & Safety untuk verifikasi manual, namun pencairan dana ke Toko Mitra TIDAK ditahan (karena kustodi barang sudah terverifikasi sah pada tahap sebelumnya) — hanya publikasi ke Beranda Publik yang ditahan hingga foto layak.

### 8.7 Konektivitas Terputus Saat Handshake atau Verifikasi
- Token QR/PIN yang sudah dibangkitkan tetap valid untuk divalidasi dalam jendela 90 detik meski perangkat sempat offline, menggunakan cache lokal terenkripsi.
- Aksi yang dilakukan offline (scan, foto) disimpan dalam antrean lokal dan disinkronkan otomatis begitu koneksi pulih, dengan timestamp asli dipertahankan untuk menjaga urutan rantai kustodi yang akurat.
- Jika sinkronisasi tidak terjadi dalam 24 jam (perangkat hilang/rusak), Admin Posko atau Trust & Safety dapat memicu proses verifikasi manual alternatif dengan bukti pendukung (foto tambahan, kesaksian pihak ketiga terverifikasi).

### 8.8 Dugaan Fraud Terbukti
- Akun yang terbukti fraud (posko fiktif, foto rekayasa, kolusi toko-kurir) langsung dinonaktifkan permanen dan seluruh dana terkait dibekukan oleh Trust & Safety (peran Approver).
- Donatur yang terdampak menerima refund penuh dari Dana Cadangan Platform, sementara kasus dapat diteruskan ke jalur hukum sesuai ketentuan yang berlaku.
- Identitas aktor yang terbukti fraud dicatat pada daftar internal (denylist) berbasis NIK/nomor rekening untuk mencegah pendaftaran ulang dengan identitas sama.

### 8.9 Donatur Mengajukan Pembatalan Setelah Dana Terkunci
- Selama dana masih berstatus "terkumpul, order belum di-broadcast": pembatalan/refund dapat diproses instan oleh sistem tanpa eskalasi.
- Setelah order di-broadcast ke Toko Mitra: pembatalan tidak lagi otomatis dan memerlukan peninjauan Trust & Safety, karena toko mungkin sudah mulai menyiapkan barang berdasarkan komitmen dana tersebut.

---

## 9. Kebutuhan Non-Fungsional (Non-Functional Requirements)

### 9.1 Keamanan & Privasi Data
- Seluruh data dalam transit dienkripsi via TLS 1.3; data sensitif at-rest (NIK, nomor rekening, foto KTP) dienkripsi menggunakan AES-256 dengan kunci terkelola secara terpisah dari database utama (envelope encryption).
- Kepatuhan terhadap UU Pelindungan Data Pribadi No. 27 Tahun 2022: consent eksplisit saat pendaftaran, hak akses/hapus data pengguna (data subject rights), serta pencatatan Data Protection Impact Assessment untuk pemrosesan data KYC skala besar.
- Foto penyintas pada Beranda Transparansi Publik tidak menampilkan nama lengkap individu; hanya inisial dan nama posko/wilayah, dengan opsi buram wajah otomatis (face-blur) dapat diaktifkan Admin Posko untuk melindungi privasi kelompok rentan (anak-anak, korban kekerasan).
- Audit log tidak dapat dihapus (append-only) untuk seluruh tindakan Trust & Safety, transaksi dana, dan perubahan status ledger.
- Rate limiting pada seluruh endpoint publik untuk mencegah scraping data posko/donatur dan serangan brute-force pada validasi PIN.
- Uji penetrasi (penetration testing) pihak ketiga wajib dilakukan sebelum peluncuran publik dan setiap 6 bulan setelahnya.

### 9.2 Kinerja (Performance)

| Skenario | Target Latensi (p95) | Catatan |
|---|---|---|
| Muat awal Peta Bencana Interaktif | < 2 detik | Termasuk clustering marker hingga 500 posko aktif |
| Broadcast order ke Toko Mitra | < 3 detik dari dana terkunci | Mencakup query geospasial + push notification |
| Validasi handshake QR/PIN | < 1 detik | Termasuk validasi HMAC dan penulisan ledger |
| Unggah & kompresi foto verifikasi | < 5 detik pada jaringan 3G | Kompresi klien 80% sebelum transmisi |
| Kalkulasi debit saldo proporsional | < 500 ms untuk 10.000 donatur bersamaan | Memanfaatkan konkurensi Goroutine di backend |

### 9.3 Skalabilitas & Keandalan
- Arsitektur backend stateless agar dapat di-scale horizontal secara otomatis saat lonjakan trafik terjadi bersamaan dengan bencana skala nasional.
- Target uptime layanan inti: 99,5% pada kondisi normal, ditingkatkan menjadi 99,9% dengan kapasitas cadangan (standby capacity) yang diaktifkan otomatis saat status siaga darurat nasional dideklarasikan BNPB.
- Database utama menggunakan replikasi multi-zona (multi-AZ) dengan Recovery Point Objective (RPO) ≤ 5 menit dan Recovery Time Objective (RTO) ≤ 30 menit.
- Circuit breaker pada integrasi payment gateway agar kegagalan pihak ketiga tidak menjatuhkan keseluruhan sistem (graceful degradation: order matching tetap berjalan meski disbursement tertunda antre).

### 9.4 Kepatuhan Regulasi
- Izin Pengumpulan Uang atau Barang (PUB) dari Kementerian Sosial RI sebagai payung hukum penggalangan dana publik untuk kebencanaan.
- Status sebagai Penyelenggara Jasa Pembayaran atau kerja sama resmi dengan PJP berizin (di bawah pengawasan Bank Indonesia) untuk aktivitas penerimaan dan penyaluran dana pihak ketiga.
- Laporan keuangan penggalangan dana kebencanaan disusun sesuai pedoman transparansi lembaga penggalang dana yang diawasi Kemensos, dipublikasikan berkala (bulanan) di luar Beranda Transparansi per-transaksi.
- Ketentuan anti pencucian uang (APU) dan pencegahan pendanaan terorisme (PPT) diterapkan pada proses KYC Toko Mitra dan pemantauan transaksi bernilai tidak wajar.

### 9.5 Aksesibilitas & Ketahanan Lapangan
- Antarmuka PWA dioptimalkan untuk perangkat kelas menengah-bawah (RAM 2–3GB) mengingat banyak relawan dan Toko Mitra di daerah menggunakan perangkat non-flagship.
- Kontras warna dan ukuran tombol mengikuti standar WCAG 2.1 AA agar tetap dapat dioperasikan dalam kondisi pencahayaan buruk atau oleh pengguna dengan keterbatasan penglihatan.
- Seluruh alur kritis (broadcast order, handshake, verifikasi akhir) berfungsi penuh dalam mode data seluler 2G/3G dengan fallback SMS untuk notifikasi (Bagian 7.11).

---

## 10. Arsitektur Teknis & Pipeline Sistem

Arsitektur mempertahankan prinsip dasar dari konsep awal — desentralisasi logika melalui Mock-Ledger bergaya blockchain dan sentralisasi kecepatan melalui backend Golang — dengan penambahan komponen yang diperlukan untuk operasional produksi: caching/geospasial cepat, antrean pesan untuk broadcast order, penyimpanan objek untuk foto, dan gateway SMS untuk ketahanan sinyal rendah.

### 10.1 Tumpukan Teknologi (Tech Stack Lengkap)

| Layer | Teknologi | Rasionalisasi |
|---|---|---|
| Frontend | Next.js + Tailwind CSS (PWA) | Server-Side Rendering menjamin pemuatan cepat; kapabilitas PWA menjaga pemindai QR & modul kamera tetap responsif di area minim sinyal |
| Backend | Golang + GORM | Goroutine memberi konkurensi tinggi untuk kalkulasi debit saldo proporsional puluhan ribu donatur secara serentak |
| Database Utama | MariaDB (Spatial + JSON) | ST_Distance_Sphere untuk geofencing; JSON untuk menyimpan log hash berantai (Merklized Mock-Ledger) |
| Cache & Geospasial Cepat | Redis (dengan modul geospasial) | Query radius toko terdekat < 50ms; penyimpanan sesi token QR/PIN dinamis dengan TTL otomatis 30 detik |
| Antrean Pesan | NATS / RabbitMQ | Menjamin broadcast order terkirim ke seluruh Toko Mitra kandidat tanpa hilang meski terjadi lonjakan trafik bencana besar |
| Penyimpanan Objek | S3-compatible Object Storage (mis. MinIO/Cloud Storage) | Menyimpan foto posko, foto verifikasi akhir, dan dokumen KYC secara terpisah dari database transaksional |
| Gateway Pembayaran | Xendit / Midtrans | Mendukung inbound QRIS/VA dan outbound disbursement otomatis ke rekening Toko Mitra |
| Notifikasi | Push Notification (FCM/APNs) + SMS Gateway (mis. Zenziva/Twilio) | Fallback SMS untuk area dengan sinyal data lemah namun sinyal seluler dasar tersedia |
| Real-time Map Update | WebSocket / Server-Sent Events | Memperbarui status posko dan lokasi kurir di peta tanpa polling berlebihan |
| Observability | Prometheus + Grafana, terintegrasi dengan sistem alerting on-call | Memantau SLA broadcast order, latensi handshake, dan tingkat kegagalan geofencing secara real-time |

### 10.2 Diagram Alur Data (Ringkas)

**Akuisisi Data:** Data ditangkap murni melalui interaksi perangkat keras pada peramban seluler. Saat serah terima logistik, sensor kamera menangkap string QR Code Dinamis; pada verifikasi akhir, sistem memaksa antarmuka kamera bawaan (mengunci akses galeri) sekaligus mengekstraksi koordinat Lintang/Bujur dan timestamp UTC melalui HTML5 Geolocation API.

**Pra-Pemrosesan & Transmisi:** Di sisi klien, gambar dikompresi hingga 80% untuk menghemat bandwidth. Data sensor, string QR, dan gambar dibungkus dalam JSON, dienkripsi, dan dikirim ke backend Golang via HTTPS asinkron; jika koneksi terputus, payload disimpan dalam antrean lokal terenkripsi (Bagian 7.11).

**Analitik & Aksi:** Backend memvalidasi geofencing, mengambil hash transaksi sebelumnya dari MariaDB, menggabungkannya dengan data baru, menghasilkan hash SHA-256 baru yang dikunci ke tabel log. Setelah integritas rantai tervalidasi, Golang memicu webhook ke payment gateway untuk eksekusi transfer fiat real-time dan mendistribusikan poin loyalitas.

### 10.3 Model Data Inti (Skema Konseptual)

Tabel berikut merangkum entitas utama beserta relasi kunci untuk memandu perancangan skema database secara rinci pada fase engineering.

| Entitas | Field Kunci | Relasi |
|---|---|---|
| users | id, role (donor/store/courier/posko_admin/trust_safety), phone, email, kyc_status, reputation_score, created_at | 1:N ke seluruh entitas transaksional sesuai peran |
| posko | id, admin_id, lat, lng, geofence_radius_m, status, disaster_type, verification_status | 1:N ke requests, N:1 ke users(admin) |
| requests | id, posko_id, item_name, unit, qty_needed, est_unit_price, funding_status | N:1 ke posko, 1:1 ke orders |
| donations | id, donor_id, request_id, amount, status(locked/refunded/used), locked_at | N:1 ke users(donor), N:1 ke requests |
| orders | id, request_id, store_id, courier_id, status, accepted_at, broadcast_radius_km | N:1 ke requests, N:1 ke users(store/courier) |
| custody_log | id, order_id, from_actor_id, to_actor_id, method(qr/pin), lat, lng, hash, prev_hash, timestamp | N:1 ke orders — rantai Merkle per order |
| disbursements | id, order_id, store_id, amount, idempotency_key, status, gateway_ref | 1:1 ke orders |
| loyalty_points | id, user_id, order_id, points, source, expires_at | N:1 ke users, N:1 ke orders |
| disputes | id, order_id, raised_by, reason, status, resolution_note, reviewer_id, approver_id | N:1 ke orders, N:1 ke users |
| flags | id, order_id/posko_id, reporter_id (nullable utk publik anonim), reason, created_at | N:1 ke orders/posko |
| audit_log | id, actor_id, action, target_type, target_id, metadata_json, timestamp | Append-only, lintas seluruh entitas |

*Catatan: skema di atas adalah rancangan konseptual tingkat PRD; normalisasi penuh, indeks, dan tipe data presisi didefinisikan pada dokumen desain teknis (Technical Design Doc) terpisah oleh tim Engineering.*

### 10.4 Spesifikasi API Inti (Ringkasan)

| Method | Endpoint | Deskripsi | Auth |
|---|---|---|---|
| POST | /api/auth/register | Registrasi aktor baru (role-based, memicu alur KYC bila perlu) | Publik |
| POST | /api/auth/otp/verify | Verifikasi OTP pendaftaran/login | Publik |
| POST | /api/posko | Membuat posko baru beserta daftar kebutuhan | Posko Admin |
| GET | /api/map/posko?bbox= | Mengambil posko dalam viewport peta (dengan clustering) | Publik |
| POST | /api/donors/deposit | Top-up saldo & alokasi ke request tertentu | Donor |
| POST | /api/requests/{id}/fund-check | Trigger internal saat dana request mencukupi target | Sistem |
| POST | /api/orders/{id}/accept | Toko menyetujui order (atomic compare-and-swap) | Store |
| POST | /api/custody/handshake | Submit hasil scan QR/PIN untuk perpindahan kustodi | Store/Courier/Posko |
| POST | /api/orders/{id}/complete | Submit foto verifikasi akhir + GPS untuk penyelesaian order | Posko Admin |
| POST | /v1/disbursements | Trigger pencairan dana ke rekening Toko Mitra | Sistem → Payment Gateway |
| GET | /api/donors/{id}/ledger | Riwayat transaksi & alokasi dana donatur | Donor (own data) |
| POST | /api/flags | Mengajukan laporan/flag atas posko atau order | Publik/Terdaftar |
| POST | /api/disputes | Mengajukan sengketa formal atas suatu order | Store/Courier/Donor |
| PATCH | /api/admin/disputes/{id}/resolve | Trust & Safety menyelesaikan sengketa (perlu peran Approver) | Trust & Safety |
| POST | /api/admin/accounts/{id}/freeze | Membekukan akun & dana terkait | Trust & Safety (Approver) |

Seluruh endpoint yang menerima data finansial atau mengubah status kustodi wajib menyertakan idempotency key pada header request untuk mencegah duplikasi akibat retry jaringan yang umum terjadi pada konektivitas area bencana.

---

## 11. Rencana Analitik & Instrumentasi

Setiap event kunci diinstrumentasi agar tim produk dapat memantau funnel dari trigger emosional di peta hingga penyelesaian pesanan, serta mendeteksi titik gesekan (friction point) secara data-driven.

| Event | Properti Kunci | Tujuan Pengukuran |
|---|---|---|
| posko_created | posko_id, disaster_type, item_count, target_amount | Volume & pola kebutuhan logistik per jenis bencana |
| map_pin_clicked | posko_id, funding_ratio, user_type | Efektivitas gradasi warna dalam memicu klik donatur |
| donation_completed | amount, request_id, is_autonomous_mode | Konversi funnel donasi & adopsi mode donasi otonom |
| order_broadcast | request_id, initial_radius_km | Baseline kepadatan Toko Mitra per wilayah |
| order_accepted | order_id, store_id, time_to_accept_sec | Kecepatan respons Toko Mitra, indikator likuiditas pasar |
| custody_handshake | order_id, method(qr/pin), gps_delta_m | Tingkat penggunaan fallback PIN & anomali GPS |
| order_completed | order_id, photo_retry_count, geofence_status | Kualitas verifikasi akhir & tingkat penolakan foto |
| dispute_raised / dispute_resolved | order_id, reason, resolution_time_hr | Kesehatan operasional & beban kerja Trust & Safety |
| flag_submitted | target_type, reason, reporter_type | Efektivitas moderasi berbasis komunitas |
| loyalty_redeemed | user_id, points, reward_type | Efektivitas program retensi |

Seluruh event di atas dialirkan ke data warehouse melalui pipeline event streaming untuk mendukung dashboard operasional real-time (Trust & Safety) maupun analisis produk berkala (mingguan/bulanan oleh tim Produk).

---

## 12. Risiko & Mitigasi

| Risiko | Dampak | Probabilitas | Mitigasi |
|---|---|---|---|
| Kolusi Toko Mitra & Kurir untuk menggelembungkan harga barang | Tinggi | Sedang | Katalog harga referensi pasar per wilayah + deteksi anomali harga otomatis + audit sampling berkala oleh Trust & Safety |
| Posko fiktif untuk pencairan dana palsu | Tinggi | Rendah–Sedang | Verifikasi Admin Posko berlapis (afiliasi organisasi/vouching), validasi geofencing ketat, deteksi duplikasi posko |
| Relawan Kurir mengalami kecelakaan/cedera saat bertugas | Tinggi | Sedang | Fase MVP: waiver risiko eksplisit + rekomendasi asuransi mandiri; Fase 2: eksplorasi kemitraan asuransi mikro kerelawanan |
| Kegagalan payment gateway saat lonjakan trafik bencana nasional | Tinggi | Sedang | Circuit breaker + antrean disbursement dengan retry otomatis + kemitraan dual-gateway (Xendit & Midtrans) untuk redundansi |
| Kebocoran data pribadi KYC (NIK, foto KTP, rekening bank) | Sangat Tinggi | Rendah | Enkripsi at-rest, akses berbasis peran dengan audit log, penetration testing berkala (Bagian 9.1) |
| Donatur kehilangan kepercayaan akibat foto/laporan yang terlihat direkayasa | Tinggi | Sedang | Forced camera + geofencing + pelaporan komunitas + eskalasi cepat Trust & Safety (Bagian 8.6) |
| Kelangkaan Toko Mitra di wilayah terpencil (luar Jawa) | Sedang | Tinggi | Kemitraan awal dengan grosir/ritel nasional berjejaring luas sebagai toko cadangan (fallback tier) |
| Ketergantungan sinyal internet di lokasi bencana | Tinggi | Tinggi | PWA offline-tolerant, cache token 90 detik, fallback SMS (Bagian 7.11) |
| Perubahan regulasi penggalangan dana/pembayaran digital | Sedang | Rendah–Sedang | Pemantauan regulasi berkelanjutan bersama penasihat hukum, desain modular pada lapisan kepatuhan agar mudah disesuaikan |
| Beban kerja Trust & Safety melampaui kapasitas saat bencana besar serentak | Sedang | Sedang | SLA berjenjang berdasarkan urgensi, kemungkinan menambah tim on-call cadangan saat status siaga darurat nasional |

---

## 13. Prioritisasi Fitur (Kerangka RICE)

Setiap fitur dinilai dengan formula RICE = (Reach × Impact × Confidence) ÷ Effort. Reach diestimasi dalam jumlah aktor terdampak per kuartal pertama pasca-peluncuran; Impact menggunakan skala massive(3)/high(2)/medium(1)/low(0.5); Confidence dalam persentase keyakinan estimasi; Effort dalam person-month.

| Fitur | Reach | Impact | Confidence | Effort (PM) | Skor RICE |
|---|---|---|---|---|---|
| Peta Interaktif + Donasi Dasar | 20.000 | Massive | 80% | 3 | 16.000 |
| Order Matching & Bidding Toko | 3.000 | High | 80% | 3 | 1.600 |
| Rantai Kustodi QR + Fallback PIN | 3.000 | High | 70% | 4 | 1.050 |
| Forced Camera & Geofencing | 3.000 | High | 80% | 2 | 2.400 |
| Beranda Transparansi Publik | 20.000 | Medium | 90% | 2 | 4.500 |
| Backoffice Trust & Safety (Modul Baru) | 500 | Massive | 90% | 3 | 450 |
| Modul Dispute & Refund (Modul Baru) | 500 | High | 70% | 3 | 233 |
| Notifikasi SMS Fallback (Modul Baru) | 3.000 | Medium | 70% | 1 | 2.100 |
| Program Loyalitas & Reward | 15.000 | Medium | 70% | 2 | 5.250 |
| Mode Offline Penuh (Fase 2) | 3.000 | High | 50% | 5 | 600 |
| Prediksi Kebutuhan Logistik (Fase 3) | 500 | Medium | 40% | 5 | 40 |

**Interpretasi:** meski skor mentah tertinggi didominasi fitur front-facing (Peta, Beranda Publik, Loyalitas), Modul Backoffice Trust & Safety dan Dispute & Refund tetap wajib masuk MVP walau skor lebih rendah — keduanya adalah prasyarat kepatuhan dan mitigasi risiko yang tidak dapat ditunda meski secara murni RICE terlihat kurang mendesak dibanding fitur pertumbuhan pengguna.

---

## 14. Roadmap & Fase Peluncuran

| Fase | Durasi | Cakupan Utama |
|---|---|---|
| 0. Purwarupa Hackathon | 30 jam | Peta interaktif, donasi dasar, order matching, QR kustodi dua tahap, forced camera — dibuktikan pada GarudaHack 7.0 (baseline dokumen ini) |
| 1. Closed Beta | 6 minggu | Tambahkan KYC dasar, backoffice Trust & Safety minimum, dispute & refund dasar, notifikasi SMS fallback; uji dengan 3–5 posko mitra terkurasi di satu wilayah |
| 2. Public MVP Launch | 8 minggu setelah beta | Seluruh modul Bagian 7.1–7.12 lengkap, kepatuhan regulasi dasar (PUB Kemensos, kerja sama PJP), target 50 posko dalam 6 bulan |
| 3. Ekspansi & Ketahanan | 3–6 bulan pasca-launch | Mode offline penuh, program reputasi lanjutan, integrasi BPBD/BNPB, multi-bahasa daerah |
| 4. Skala & Prediksi | 6–12 bulan pasca-launch | Marketplace CSR korporasi, model prediktif kebutuhan logistik, ekspansi ke bantuan non-logistik |

---

## 15. Asumsi & Pertanyaan Terbuka

### 15.1 Asumsi Kunci
- Toko Mitra memiliki akses rekening bank yang dapat diintegrasikan dengan payment gateway (belum mencakup toko yang sepenuhnya cash-only tanpa rekening formal).
- Estimasi harga referensi katalog barang dapat diperbarui berkala dan cukup akurat untuk mendeteksi anomali harga tanpa memerlukan negosiasi manual setiap pesanan.
- Ketersediaan sinyal seluler dasar (2G) tetap ada di sebagian besar lokasi bencana meski sinyal data terganggu, sebagai prasyarat fallback SMS.

### 15.2 Pertanyaan Terbuka untuk Pemangku Kepentingan
1. Apakah PijarNusa akan beroperasi sebagai entitas penggalang dana independen dengan izin PUB sendiri, atau bermitra dengan lembaga kemanusiaan berizin yang sudah ada?
2. Siapa yang menanggung Dana Cadangan Platform untuk skenario barang hilang/fraud (Bagian 8.5, 8.8) — margin platform, dana talangan investor, atau skema asuransi pihak ketiga?
3. Apakah dibutuhkan verifikasi lapangan fisik (mis. kerja sama dengan BPBD daerah) untuk posko dengan nilai pendanaan besar (di atas ambang tertentu) sebagai lapisan kepercayaan tambahan?
4. Bagaimana kebijakan data penyintas anak-anak pada foto verifikasi akhir — apakah face-blur otomatis wajib atau opsional (Bagian 9.1)?
5. Berapa besar alokasi awal Dana Cadangan Platform yang perlu disiapkan sebelum peluncuran publik untuk menyerap potensi klaim sengketa di bulan-bulan pertama?

---

## 16. Glosarium

| Istilah | Definisi |
|---|---|
| Mock-Ledger | Struktur pencatatan transaksi berantai (mirip blockchain) yang disimpan di MariaDB menggunakan hash SHA-256 berjenjang untuk menjamin data tidak dapat diubah retroaktif tanpa merusak seluruh rantai |
| Chain of Custody | Rantai kustodi — catatan berjenjang siapa memegang barang secara sah di setiap titik waktu, dari Toko hingga Posko |
| Forced Camera | Mekanisme yang mengunci akses galeri perangkat dan memaksa penggunaan kamera langsung untuk mencegah unggahan foto lama/rekayasa |
| Geofencing | Validasi berbasis koordinat GPS untuk memastikan suatu aksi (mis. verifikasi akhir) dilakukan di lokasi yang sah |
| Overfunding / Underfunding | Kondisi dana terkumpul melebihi atau tidak mencapai kebutuhan tercatat suatu posko |
| Trust & Safety | Tim/peran internal yang memoderasi platform, menangani sengketa, dan menjaga integritas dana serta identitas pengguna |
| RICE | Kerangka prioritisasi fitur: Reach × Impact × Confidence ÷ Effort |
| PUB | Izin Pengumpulan Uang atau Barang dari Kementerian Sosial RI, payung hukum penggalangan dana publik di Indonesia |
| PJP | Penyelenggara Jasa Pembayaran, entitas berizin Bank Indonesia yang mengelola aktivitas pembayaran digital |

### Lampiran

Activity Diagram (Alur Pemindaian QR & Rantai Kustodi Logistik) dan Sequence Diagram (Urutan Interaksi Temporal Sistem) dari dokumen konsep awal GarudaHack 7.0 tetap berlaku sebagai referensi visual teknis dan dilampirkan pada repositori desain terpisah; PRD ini menambahkan detail perilaku sistem di atasnya melalui Bagian 7 dan 8.
