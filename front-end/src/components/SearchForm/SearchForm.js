'use client'


import {useEffect, useState} from "react";
import axios from "axios";
import Card from "@/components/Card/Card";

const SearchForm = () => {
    const [search, setSearch] = useState('');
    const [products, setProducts] = useState([]);
    const handleSearch = async (e) => {
        e.preventDefault()
        if (search === '') {
            return
        }
        await axios.get('http://localhost:8088/product/search?query=' + search)
            .then((response) => {
                setProducts(response.data.products)
            }).catch((error) => {
                console.log(error)
            })
    }
    useEffect(() => {

    }, []);
    return (
        <div>
            <form onSubmit={handleSearch} className="mt-12 max-w-md mx-auto">
                <label htmlFor="default-search"
                       className="mb-2 text-sm font-medium text-gray-900 sr-only dark:text-white">Search</label>
                <div className="relative">
                    <div className="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
                        <svg className="w-4 h-4 text-gray-500 dark:text-gray-400" aria-hidden="true"
                             xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 20 20">
                            <path stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                                  d="m19 19-4-4m0-7A7 7 0 1 1 1 8a7 7 0 0 1 14 0Z"/>
                        </svg>
                    </div>
                    <input value={search} onChange={(e) => setSearch(e.target.value)} type="search" id="default-search"
                           className="block w-full p-4 ps-10 text-sm text-gray-900 border border-gray-300 rounded-lg bg-gray-50 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
                           placeholder="Search Mockups, Logos..." required/>
                    <button onClick={handleSearch} type="submit"
                            className="text-white absolute end-2.5 bottom-2.5 bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-4 py-2 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">Search
                    </button>
                </div>
            </form>
            <div className="p-8 mt-12 flex items-center justify-between flex-wrap" style={{rowGap:'16px',columnGap:'8px'}}>
                {
                    products?.length && products.map(item => (
                        <Card current_price={item.current_price} title={item.name} key={item.id} img={item.images[0].address}/>
                    ))
                }
            </div>
        </div>

    )
}
export default SearchForm