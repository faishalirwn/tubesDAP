/*
Tubes DAP Topik 14: Vending Machine
Anggota kelompok:
- Muhamad Elang Ramadhan
- Muhamad Faishal Irawan
- Mochammad Iqbal
*/
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

// Makanan adalah representasi dari makanan ringan dalam vending machine
type Makanan struct {
	nama                   string
	harga, jumlah, terjual int
}

var (
	listMakanan               []Makanan
	user                      string
	pilihan, saldo, kembalian int
	mati                      bool = false
)

func clearScr() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func bacaFile(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	lines := strings.Split(string(content), "\n")
	for i := 0; i < len(lines); i++ {
		var makanan Makanan
		splitMakanan := strings.Split(string(lines[i]), ", ")
		makanan.nama = splitMakanan[0]
		makanan.harga, _ = strconv.Atoi(splitMakanan[1])
		makanan.jumlah, _ = strconv.Atoi(splitMakanan[2])
		if makanan.jumlah > 15 {
			makanan.jumlah = 15
		}
		if len(splitMakanan) > 3 {
			makanan.terjual, _ = strconv.Atoi(splitMakanan[3])
		}
		listMakanan = append(listMakanan, makanan)
	}
}

func pilihMakanan() {
	clearScr()
	w := tabwriter.NewWriter(os.Stdout, 5, 5, 20, ' ', 0)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", "No.", "Nama", "Harga", "Stok")
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", "-----", "-----", "-----", "-----")
	for i := 0; i < len(listMakanan); i++ {
		fmt.Fprintf(w, "%d\t%s\t%d\t%d\t\n ", i+1, listMakanan[i].nama, listMakanan[i].harga, listMakanan[i].jumlah)
	}
	w.Flush()

	fmt.Println("Pilih nomor makanan yang Anda inginkan:")
	fmt.Scanln(&pilihan)
	if user == "1" {
		for listMakanan[pilihan-1].jumlah == 0 {
			fmt.Println("Maaf stok makanan habis, silakan pilih nomor makanan lain:")
			fmt.Scanln(&pilihan)
		}
	}
}

// Mendeteksi uang, menambahnya ke saldo dan mengatur kembalian, lalu membeli makanan
func beliMakanan(pilihan int) {
	var (
		uang  int
		harga int = listMakanan[pilihan].harga
	)
	clearScr()
	fmt.Println("===============")
	fmt.Printf("%v\nHarga: %v\n", listMakanan[pilihan].nama, listMakanan[pilihan].harga)
	fmt.Println("===============")
	fmt.Print("Masukkan uang: ")
	fmt.Scanln(&uang)
	for saldo < harga {
		for uang != 1000 && uang != 2000 && uang != 5000 && uang != 10000 {
			clearScr()
			fmt.Printf("Saldo: %v\nPecahan uang tidak valid. Hanya menerima pecahan 1000, 2000, 5000, dan 10000.\nButuh Rp. %v lagi untuk membeli. Masukkan uang:\n", saldo, harga-saldo)
			fmt.Scanln(&uang)
		}
		saldo += uang
		if saldo < harga {
			clearScr()
			fmt.Printf("Saldo: %v\nButuh Rp. %v lagi untuk membeli. Masukkan uang:\n", saldo, harga-saldo)
			fmt.Scanln(&uang)
		} else {
			clearScr()
			kembalian = saldo - harga
			fmt.Printf("Saldo: %v\nSaldo mencukupi. Kembalian: %v\n", saldo, kembalian)
		}
	}
	fmt.Println(listMakanan[pilihan].nama, "berhasil dibeli. Terima kasih, selamat menikmati.")
	listMakanan[pilihan].jumlah--
	listMakanan[pilihan].terjual++
}

func cekMakanan(pilihan int) {
	clearScr()
	fmt.Println(listMakanan[pilihan].nama)
	fmt.Printf("Terjual: %v, Stok: %v\n", listMakanan[pilihan].terjual, listMakanan[pilihan].jumlah)
}

func cekStokMakanan(pilihan int) {
	clearScr()
	if pilihan == 1 {
		sort.Slice(listMakanan, func(p, q int) bool {
			return listMakanan[p].jumlah < listMakanan[q].jumlah
		})
	} else if pilihan == 2 {
		sort.Slice(listMakanan, func(p, q int) bool {
			return listMakanan[p].jumlah > listMakanan[q].jumlah
		})
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', 0)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", "No", "Nama", "Harga", "Stok")
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", "-----", "-----", "-----", "-----")
	for i := 0; i < len(listMakanan); i++ {
		fmt.Fprintf(w, "%d\t%s\t%d\t%d\t\n ", i+1, listMakanan[i].nama, listMakanan[i].harga, listMakanan[i].jumlah)
	}
	w.Flush()
}

func cekPalingLaku() {
	var (
		palingLaku int
		makanan    Makanan
	)
	clearScr()
	for i := 0; i < len(listMakanan); i++ {
		if listMakanan[i].terjual >= palingLaku {
			palingLaku = listMakanan[i].terjual
			makanan = listMakanan[i]
		}
	}

	if palingLaku == 0 {
		fmt.Println("Belum ada makanan yang paling laku.")
	} else {
		fmt.Println("Makanan paling laku:")
		fmt.Printf("%v\nTerjual: %v, Stok: %v\n", makanan.nama, makanan.terjual, makanan.jumlah)
	}

}

func restock(indexMakanan int) {
	var (
		makanan              string = listMakanan[indexMakanan].nama
		stokLama             int    = listMakanan[indexMakanan].jumlah
		penambahan, stokBaru int
	)
	clearScr()
	fmt.Printf("%v\nStok: %v\n", makanan, stokLama)
	fmt.Println("Jumlah makanan yang ingin ditambahkan:")
	fmt.Scanln(&penambahan)
	stokBaru = stokLama + penambahan
	for stokBaru > 15 {
		fmt.Println("Maks stok: 15. Silakan masukkan kembali jumlah makanan yang ingin ditambahkan:")
		fmt.Scanln(&penambahan)
		stokBaru = stokLama + penambahan
	}
	listMakanan[indexMakanan].jumlah = stokBaru
	fmt.Println("Makanan berhasil di restock.\n")
	fmt.Printf("%v\nStok: %v", makanan, stokBaru)
}

func tulisFile() {
	var text strings.Builder
	for i := 0; i < len(listMakanan); i++ {
		if i == len(listMakanan)-1 {
			fmt.Fprintf(&text, "%v, %v, %v, %v", listMakanan[i].nama, listMakanan[i].harga, listMakanan[i].jumlah, listMakanan[i].terjual)
		} else {
			fmt.Fprintf(&text, "%v, %v, %v, %v\n", listMakanan[i].nama, listMakanan[i].harga, listMakanan[i].jumlah, listMakanan[i].terjual)
		}

	}
	textByte := []byte(text.String())
	// Sesuaikan path file data makanan
	err := ioutil.WriteFile("./data_makanan.txt", textByte, 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	for mati == false {
		// Sesuaikan path file data makanan
		bacaFile("./data_makanan.txt")

		fmt.Println("............................................................................................................................................................................................")
		fmt.Println("..####...##..##....##.....####...##..##.......##..##..######..##..##..#####...######..##..##...####........##...##...####....####...##..##..######..##..##..######......##............####..")
		fmt.Println(".##......###.##..##..##..##..##..##.##........##..##..##......###.##..##..##....##....###.##..##...........###.###..##..##..##..##..##..##....##....###.##..##.........###...........##..##.")
		fmt.Println("..####...##.###..######..##......####.........##..##..####....##.###..##..##....##....##.###..##.###.......##.#.##..######..##......######....##....##.###..####........##...........##..##.")
		fmt.Println(".....##..##..##..##..##..##..##..##.##.........####...##......##..##..##..##....##....##..##..##..##.......##...##..##..##..##..##..##..##....##....##..##..##..........##.....00....##..##.")
		fmt.Println("..####...##..##..##..##...####...##..##.........##....######..##..##..#####...######..##..##...####........##...##..##..##...####...##..##..######..##..##..######....######...00.....####..")
		fmt.Println("............................................................................................................................................................................................")

		fmt.Println("Silahkan Pilih 1 untuk membeli. Pilih 0 bagi pengelola")
		fmt.Scanln(&user)
		for user != "1" && user != "0" {
			fmt.Println("Pilihan tidak valid")
			fmt.Println("Pilih 1 untuk membeli. Pilih 0 bagi pengelola")
			fmt.Scanln(&user)
		}
		if user == "1" {
			pilihMakanan()
			beliMakanan(pilihan - 1)
			tulisFile()
		} else if user == "0" {
			var pilihanMenu int
			clearScr()
			fmt.Println("Menu: \n1. Cek terjual dan stok makanan\n2. Stok semua makanan\n3. Makanan paling laku\n4. Restock makanan\n99. Matikan mesin")
			fmt.Scanln(&pilihanMenu)
			switch pilihanMenu {
			case 1:
				pilihMakanan()
				cekMakanan(pilihan - 1)
			case 2:
				var pilihanOpsi int
				fmt.Println("Pilih opsi:\n1.Ascending\n2.Descending")
				fmt.Scanln(&pilihanOpsi)
				cekStokMakanan(pilihanOpsi)
			case 3:
				cekPalingLaku()
			case 4:
				pilihMakanan()
				restock(pilihan - 1)
				tulisFile()
			case 99:
				mati = true
			}
		}
		saldo = 0
		listMakanan = nil
		time.Sleep(4 * time.Second)
		clearScr()
	}

}
